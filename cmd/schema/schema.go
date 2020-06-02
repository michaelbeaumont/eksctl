/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Modifications:
  Copyright 2020 Weaveworks
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	version7  = "http://json-schema-org/draft-07/schema#"
	defPrefix = "#/definitions/"
)

var (
	regexpDefaults = regexp.MustCompile("(.*)Defaults to `(.*)`")
	regexpExample  = regexp.MustCompile("(.*)For example: `(.*)`")
	pTags          = regexp.MustCompile("(<p>)|(</p>)")

	// patterns for enum-type values
	enumValuePattern     = "^[ \t]*`(?P<name>[^`]+)`([ \t]*\\(default\\))?: .*$"
	regexpEnumDefinition = regexp.MustCompile("(?m).*Valid [a-z]+ are((\\n" + enumValuePattern + ")*)")
	regexpEnumValues     = regexp.MustCompile("(?m)" + enumValuePattern)
)

// GenerateSchema is the entrypoint for schema generation
func GenerateSchema(root, version, rootRef string) (Schema, error) {
	strict := false

	input := filepath.Join(root, "pkg", "apis", "eksctl.io", version)

	generator := schemaGenerator{
		strict: strict,
	}

	schema, err := generator.generateRoot(input, version, rootRef)
	if err != nil {
		return Schema{}, errors.Wrapf(err, "unable to generate schema for version %q", version)
	}
	return schema, nil

}

func fieldTag(field *ast.Field) reflect.StructTag {
	if field.Tag == nil {
		return ""
	}
	tag := strings.Replace(field.Tag.Value, "`", "", -1)
	return reflect.StructTag(tag)
}

func jsonPropName(tag reflect.StructTag) string {
	jsonField := tag.Get("json")

	return strings.Split(jsonField, ",")[0]
}

//nolint:golint,goconst
func setTypeOrRef(def *Definition, typeName string) {
	switch typeName {
	case "string":
		def.Type = "string"
	case "bool":
		def.Type = "boolean"
	case "int", "int64", "int32":
		def.Type = "integer"
	case "float64":
		def.Type = "number"
	case "byte":
		def.Type = "string" // TODO mediaEncoding
	default:
		def.Ref = defPrefix + typeName
	}
}

// findTypeSpecFromIdent takes a SelectorExpr "pkg.Thing" and gets a TypeSpec if
// it can
func findTypeSpecFromIdent(it *ast.SelectorExpr, imports map[string]*ast.Object) (string, *ast.TypeSpec, error) {
	importSpec := it.X.(*ast.Ident).Obj.Decl.(*ast.ImportSpec)
	importPath := importSpec.Path.Value
	importedPkg, err := DoImport(imports, importPath[1:len(importPath)-1])
	if err != nil {
		return "", nil, errors.Wrapf(err, "couldn't handle struct field")
	}
	scope := importedPkg.Data.(*ast.Scope)
	typeSpec := scope.Objects[it.Sel.Name].Decl.(*ast.TypeSpec)
	inlineName := fmt.Sprintf("%s.%s", it.X.(*ast.Ident).Name, it.Sel.Name)
	return inlineName, typeSpec, nil
}

func (g *schemaGenerator) newDefinition(
	name string, t ast.Expr, comment string, tag reflect.StructTag, definitions map[string]*Definition, imports map[string]*ast.Object,
) *Definition {
	def := &Definition{
		tags: string(tag),
	}

	switch tt := t.(type) {
	case *ast.Ident:
		typeName := tt.Name
		setTypeOrRef(def, typeName)

		switch typeName {
		case "bool":
			def.Default = "false"
		}

	case *ast.StarExpr:
		if ident, ok := tt.X.(*ast.Ident); ok {
			typeName := ident.Name
			setTypeOrRef(def, typeName)
		}

	case *ast.ArrayType:
		def.Type = "array"
		def.Items = g.newDefinition("", tt.Elt, "", "", definitions, imports)
		if def.Items.Ref == "" {
			// TODO when exactly do we have default []?
		}

	case *ast.MapType:
		def.Type = "object"
		def.Default = "{}"
		def.AdditionalProperties = g.newDefinition("", tt.Value, "", "", definitions, imports)

	case *ast.StructType:
		for _, field := range tt.Fields.List {
			tag := fieldTag(field)
			fieldName := jsonPropName(tag)

			if len(field.Names) == 0 {
				var inlineName string
				switch it := field.Type.(type) {
				case *ast.Ident:
					inlineName = it.Name
				case *ast.StarExpr:
					if iit, ok := it.X.(*ast.Ident); ok {
						inlineName = iit.Name
					}
				case *ast.SelectorExpr:
					var typeSpec *ast.TypeSpec
					var err error
					inlineName, typeSpec, err = findTypeSpecFromIdent(it, imports)
					if err != nil {
						panic(errors.Wrapf(err, "Couldn't import type from identifier"))
					}
					definitions[inlineName] = g.newDefinition(inlineName, typeSpec.Type, typeSpec.Doc.Text(), "", definitions, imports)
				default:
					panic(errors.Errorf("Unexpected inline field type %v", it))
				}
				def.PreferredOrder = append(def.PreferredOrder, "<inline>")
				def.inlines = append(def.inlines, &Definition{
					Ref: defPrefix + inlineName,
				})
				continue
			}

			if fieldName == "" {
				continue
			}

			if strings.Contains(string(tag), "required") {
				def.Required = append(def.Required, fieldName)
			}

			if def.Properties == nil {
				def.Properties = make(map[string]*Definition)
			}

			def.PreferredOrder = append(def.PreferredOrder, fieldName)
			// TODO handle imported types here
			def.Properties[fieldName] = g.newDefinition(field.Names[0].Name, field.Type, field.Doc.Text(), tag, definitions, imports)
			def.AdditionalProperties = false
		}
	}

	err := HandleComment(name, comment, def, g.strict)
	if err != nil {
		panic(err)
	}

	return def
}

func (g *schemaGenerator) generateRoot(inputPath string, pkgName string, rootRef string) (Schema, error) {
	dir := ParseDir(inputPath)
	pkg, ok := dir[pkgName]
	if !ok {
		return Schema{}, errors.Errorf("Couldn't find package %s", pkgName)
	}
	schemaPkg, _ := ast.NewPackage(token.NewFileSet(), pkg.Files, DoImport, nil)
	return g.generatePackage(schemaPkg, rootRef)
}

func (g *schemaGenerator) generatePackage(schemaPkg *ast.Package, rootRef string) (Schema, error) {
	var preferredOrder []string
	definitions := make(map[string]*Definition)

	if schemaPkg.Scope == nil {
		return Schema{}, errors.Errorf("Nil scope in pkg %s", schemaPkg.Name)
	}
	// Generate definitions for types from this package
	for _, obj := range schemaPkg.Scope.Objects {
		if obj.Decl == nil {
			continue
		}
		typeSpec, ok := obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		name := typeSpec.Name.Name
		preferredOrder = append(preferredOrder, name)
		definitions[name] = g.newDefinition(name, typeSpec.Type, typeSpec.Doc.Text(), "", definitions, schemaPkg.Imports)
	}

	// Handle embedded fields
	var inlines []string

	for _, k := range preferredOrder {
		def := definitions[k]
		if len(def.inlines) == 0 {
			continue
		}

		for _, inlineStruct := range def.inlines {
			ref := strings.TrimPrefix(inlineStruct.Ref, defPrefix)
			inlines = append(inlines, ref)
		}

		inlineIndex := 0
		var defPreferredOrder []string
		for _, k := range def.PreferredOrder {
			if k != "<inline>" {
				defPreferredOrder = append(defPreferredOrder, k)
				continue
			}

			inlineStruct := def.inlines[inlineIndex]

			ref := strings.TrimPrefix(inlineStruct.Ref, defPrefix)
			inlineStructRef := definitions[ref]
			if inlineStructRef == nil {
				return Schema{}, errors.Errorf("Couldn't find struct ref %s in definitions", ref)
			}

			if def.Properties == nil {
				def.Properties = make(map[string]*Definition, len(inlineStructRef.Properties))
			}
			for k, v := range inlineStructRef.Properties {
				def.Properties[k] = v
			}
			defPreferredOrder = append(defPreferredOrder, inlineStructRef.PreferredOrder...)
			def.Required = append(def.Required, inlineStructRef.Required...)
			inlineIndex++
		}
		def.PreferredOrder = defPreferredOrder

	}

	for _, ref := range inlines {
		delete(definitions, ref)
	}

	s := Schema{
		Version: version7,
		Definition: &Definition{
			Type: "object",
			Ref:  defPrefix + rootRef,
		},
		Definitions: definitions,
	}
	if _, ok := definitions[rootRef]; !ok {
		return s, errors.Errorf("Couldn't find ref %s in definitions", rootRef)
	}
	return s, nil
}

// ToJSON serializes and makes sure HTML description are not escaped
func ToJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
