package apidoc

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// LoadBackend parse paths in a framework.Backend into apidoc paths,
// methods, etc. It will infer path and body fields from when they're
// provided, and use existing path help if available.
func LoadBackend(prefix string, backend *framework.Backend, doc *Document) {
	for _, p := range backend.Paths {
		paths := procFrameworkPath(p)
		doc.Mounts[prefix] = append(doc.Mounts[prefix], paths...)
	}
}

// procFrameworkPath parses a framework.Path into one or more apidoc.Paths.
func procFrameworkPath(p *framework.Path) []Path {
	var docPaths []Path
	var httpMethod string

	paths := expandPattern(p.Pattern)

	for _, path := range paths {
		methods := make(map[string]Method)

		for opType := range p.Callbacks {
			m := Method{
				Summary:     cleanString(p.HelpSynopsis),
				Description: cleanString(p.HelpDescription),
			}
			switch opType {
			case logical.CreateOperation:
				httpMethod = "POST"
				m.Responses = []Response{StdRespNoContent}
			case logical.UpdateOperation:
				httpMethod = "PUT"
				m.Responses = []Response{StdRespNoContent}
			case logical.DeleteOperation:
				httpMethod = "DELETE"
				m.Responses = []Response{StdRespNoContent}
			case logical.ReadOperation, logical.ListOperation:
				httpMethod = "GET"
				m.Responses = []Response{StdRespOK}
			default:
				panic(fmt.Sprintf("unknown operation type %v", opType))
			}

			fieldSet := make(map[string]bool)
			params := pathFields(path)

			// Extract path fields
			for _, param := range params {
				fieldSet[param] = true
				typ, sub := convertType(p.Fields[param].Type)
				m.PathFields = append(m.PathFields, Property{
					Name:        param,
					Type:        typ,
					SubType:     sub,
					Description: cleanString(p.Fields[param].Description),
				})
			}

			sort.Slice(m.PathFields, func(i, j int) bool {
				return m.PathFields[i].Name < m.PathFields[j].Name
			})

			// It's assumed that any fields not present in the path can be part of
			// the body for POST/PUT methods.
			if httpMethod == "POST" || httpMethod == "PUT" {
				for name, field := range p.Fields {
					if !fieldSet[name] {
						typ, sub := convertType(field.Type)
						m.BodyFields = append(m.BodyFields, Property{
							Name:        name,
							Description: cleanString(field.Description),
							Type:        typ,
							SubType:     sub,
						})
					}
				}
			}
			sort.Slice(m.BodyFields, func(i, j int) bool {
				return m.BodyFields[i].Name < m.BodyFields[j].Name
			})

			methods[httpMethod] = m
		}

		if len(methods) > 0 {
			newPath := Path{
				Pattern: path,
				Methods: methods,
			}
			docPaths = append(docPaths, newPath)
		}
	}

	return docPaths
}

// Regexen for handling optional and named parameters
var optRe = regexp.MustCompile(`(?U)\(.*\)\?`)
var reqdRe = regexp.MustCompile(`\(\?P<(\w+)>[^)]*\)`)
var cleanRe = regexp.MustCompile("[()$?]")

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) []string {

	// This construct is added by GenericNameRegex and is much easier to remove now
	// than to compensate for in the other regexes.
	pattern = strings.Replace(pattern, `\w(([\w-.]+)?\w)?`, "", -1)

	// expand all optional regex elements into two paths. This approach is really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase beyond that anyway.
	paths := []string{pattern}

	for i := 0; i < len(paths); i++ {
		p := paths[i]
		match := optRe.FindStringIndex(p)
		if match != nil {
			paths[i] = p[0:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]
			paths = append(paths, p[0:match[0]]+p[match[1]:])
			i--
		}
	}

	// replace named parameters (?P<foo>) with {foo}
	replacedPaths := make([]string, 0)
	for _, path := range paths {
		result := reqdRe.FindAllStringSubmatch(path, -1)
		if result != nil {
			for _, p := range result {
				par := p[1]
				path = strings.Replace(path, p[0], fmt.Sprintf("{%s}", par), 1)
			}
		}
		path = cleanRe.ReplaceAllString(path, "")
		replacedPaths = append(replacedPaths, path)
	}
	return replacedPaths
}

// convertType translates a FieldType into an OpenAPI type.
// In the case of arrays, a subtype is returns as well.
func convertType(t framework.FieldType) (string, string) {
	var ret, sub string

	switch t {
	case framework.TypeString, framework.TypeNameString, framework.TypeKVPairs:
		ret = "string"
	case framework.TypeInt, framework.TypeDurationSecond:
		ret = "number"
	case framework.TypeBool:
		ret = "boolean"
	case framework.TypeMap:
		ret = "object"
	case framework.TypeSlice, framework.TypeStringSlice, framework.TypeCommaStringSlice:
		ret = "array"
		sub = "string"
	case framework.TypeCommaIntSlice:
		ret = "array"
		sub = "number"
	default:
		log.Fatalf("Unsupported type %d", t)
	}

	return ret, sub
}

// cleanString prepares s for inclusion in the output YAML. This is currently just
// basic escaping, truncating, and wrapping in quotes.
func cleanString(s string) string {
	s = strings.TrimSpace(s)

	// TODO: no truncation for now.
	//if idx := strings.Index(s, "\n"); idx != -1 {
	//	s = s[0:idx] + "..."
	//}

	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, `"`, `\"`, -1)

	return fmt.Sprintf(`"%s"`, s)
}
