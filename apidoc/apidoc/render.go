package apidoc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"text/template"
)

type OAPIRenderer struct {
	Output   io.Writer
	Template string
	Version  int
}

func NewOAPIRenderer(version int) (*OAPIRenderer, error) {
	if version != 2 {
		return nil, errors.New("Sorry, only Open API version 2 is supported!")
	}

	return &OAPIRenderer{
		Output:   os.Stdout,
		Template: Tmpl_oapi2,
		Version:  version,
	}, nil
}

func (r *OAPIRenderer) Render(doc *Document) {
	funcs := map[string]interface{}{
		"indent": funcIndent,
		"lower":  strings.ToLower,
	}

	tmpl, _ := template.New("root").Funcs(funcs).Parse(r.Template)
	tmpl.Execute(r.Output, doc)
}

func funcIndent(count int, text string) string {
	var buf bytes.Buffer
	prefix := strings.Repeat(" ", count)
	scan := bufio.NewScanner(strings.NewReader(text))
	for scan.Scan() {
		buf.WriteString(prefix + scan.Text() + "\n")
	}

	return strings.TrimRight(buf.String(), "\n")
}
