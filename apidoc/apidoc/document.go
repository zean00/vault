package apidoc

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/vault/version"
)

// Document is a set a API documentation. The structure of it and its descendants roughly
// follow the organization of OpenAPI, but it is not rigidly tied to that. Additional elements
// can be added, and many OpenAPI constructs are omitted. It is meant as an itermediate format
// from which OpenAPI or other targets can be generated.
type Document struct {
	Version string
	Mounts  map[string][]Path
}

func NewDoc() *Document {
	return &Document{
		Version: version.GetVersion().Version,
		Mounts:  make(map[string][]Path),
	}
}

func (d *Document) AddPath(mount string, p ...Path) {
	d.Mounts[mount] = append(d.Mounts[mount], p...)
}

// PathList returns a flat list of fully expanded path.
func (d *Document) PathList() []Path {
	var paths []Path
	for mount, pathe := range d.Mounts {
		for _, p := range pathe {
			p.Pattern = fmt.Sprintf("/%s/%s", mount, p.Pattern)
			paths = append(paths, p)
		}
	}

	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Pattern < paths[j].Pattern
	})

	return paths
}

// Path is the structure for a single path, including all of its methods.
// The path description is kept split as: /<prefix>/<pattern>
type Path struct {
	Pattern string
	Methods map[string]Method
}

func NewPath(pattern string) Path {
	return Path{
		Pattern: pattern,
		Methods: make(map[string]Method),
	}
}

func (p *Path) AddMethod(m Method) {
	p.Methods[m.httpMethod] = m
}

func (p *Path) Prefix() string {
	s := strings.TrimLeft(p.Pattern, "/")
	sp := strings.Split(s, "/")
	return sp[0]
}

type Method struct {
	Summary     string
	Description string
	PathFields  []Property
	BodyFields  []Property
	Responses   []Response
	httpMethod  string
}

func NewMethod(httpMethod string, summary string) Method {
	ret := Method{
		Summary:    summary,
		httpMethod: httpMethod,
	}

	switch httpMethod {
	case "GET", "PUT", "POST", "HEAD", "DELETE":
	default:
		log.Fatalf("unsupported method: %s", httpMethod)
	}

	return ret
}

// WIP
func (m *Method) AddResponse(code int, example string) {
	var description string
	switch code {
	case 200:
		description = "OK"
	case 204:
		description = "empty body"
	}
	m.Responses = append(m.Responses, NewResponse(code, description, example))
}

type Property struct {
	Name        string
	Type        string
	SubType     string
	Description string
}

func NewProperty(name, typ, description string) Property {
	p := Property{
		Name:        name,
		Description: description,
	}
	parts := strings.Split(typ, "/")
	p.Type = parts[0]
	if len(parts) > 1 && parts[0] == "array" {
		p.SubType = parts[1]
	}

	return p
}

type Response struct {
	Code        int
	Description string
	Example     string
}

func NewResponse(code int, description, example string) Response {
	example = strings.TrimSpace(example)
	example = strings.Replace(example, "\t", "  ", -1)
	return Response{
		Code:        code,
		Description: description,
		Example:     example,
	}
}

func pathFields(pattern string) []string {
	pathFieldsRe := regexp.MustCompile(`{(\w+)}`)

	r := pathFieldsRe.FindAllStringSubmatch(pattern, -1)
	ret := make([]string, 0, len(r))
	for _, t := range r {
		ret = append(ret, t[1])
	}
	return ret
}

var StdRespOK = Response{
	Code:        200,
	Description: "OK",
}

var StdRespNoContent = Response{
	Code:        204,
	Description: "empty body",
}
