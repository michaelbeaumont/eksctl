package main

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	blackfriday "github.com/russross/blackfriday/v2"
)

// HandleComment interprets as much as it can from the comment and saves this
// information in the Definition
func HandleComment(name, comment string, def *Definition, strict bool) error {
	if strict && name != "" {
		if !strings.HasPrefix(comment, name+" ") {
			return errors.Errorf("comment should start with field name on field %s", name)
		}
	}

	// process enums before stripping out newlines
	if m := regexpEnumDefinition.FindStringSubmatch(comment); m != nil {
		enums := make([]string, 0)
		if n := regexpEnumValues.FindAllStringSubmatch(m[1], -1); n != nil {
			for _, matches := range n {
				enums = append(enums, matches[1])
			}
			def.Enum = enums
		}
	}

	// Remove kubernetes-style annotations from comments
	description := strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(comment, "+required", ""),
				"+optional", "",
			), "\n", " ",
		),
	)

	// Extract default value
	if m := regexpDefaults.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		def.Default = m[2]
	}

	// Extract example
	if m := regexpExample.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		def.Examples = []string{m[2]}
	}

	// Remove type prefix
	description = regexp.MustCompile("^"+name+" (\\*.*\\* )?((is (the )?)|(are (the )?)|(lists ))?").ReplaceAllString(description, "$1")

	if strict && name != "" {
		if description == "" {
			return errors.Errorf("no description on field %s", name)
		}
		if !strings.HasSuffix(description, ".") {
			return errors.Errorf("description should end with a dot on field %s", name)
		}
	}
	def.Description = description

	// Convert to HTML
	html := string(blackfriday.Run([]byte(description), blackfriday.WithNoExtensions()))
	def.HTMLDescription = strings.TrimSpace(pTags.ReplaceAllString(html, ""))
	return nil
}
