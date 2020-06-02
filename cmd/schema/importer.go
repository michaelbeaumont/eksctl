package main

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func ignoreTestFiles(file os.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}

// ParseDir returns a map of packages
func ParseDir(path string) map[string]*ast.Package {
	dir, _ := parser.ParseDir(token.NewFileSet(), path, ignoreTestFiles, parser.ParseComments)
	return dir
}

// DoImport creates an Object for the package at `path`, caching in `imports`
func DoImport(imports map[string]*ast.Object, path string) (pkg *ast.Object, err error) {
	if imports[path] != nil {
		return imports[path], nil
	}
	// Find out where our package is
	imported, err := build.Import(path, ".", build.FindOnly)
	if err != nil {
		return nil, err
	}
	// Just take the first package from that directory
	// Ignore errors for now, there will be some
	dir := ParseDir(imported.Dir)
	for _, p := range dir {
		name := path[strings.LastIndex(path, "/")+1:]
		schemaPkg, _ := ast.NewPackage(token.NewFileSet(), p.Files, nil, nil)
		imports[path] = &ast.Object{
			Kind: ast.Pkg,
			Name: name,
			Decl: nil, // An ImportSpec "should" go here, but we don't need it
			Data: schemaPkg.Scope,
			Type: nil,
		}
		return imports[path], nil
	}
	return nil, errors.Errorf("Couldn't import package %s", path)
}
