package jstmpl

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"text/template"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Process(out io.Writer, model *Root, tmpl []byte) error {
	t := template.Must(template.New("jstmpl").Delims("/*", "*/").Funcs(map[string]interface{}{
		"notLast":               notLast,
		"spaceToUpperCamelCase": spaceToUpperCamelCase,
		"snakeToLowerCamelCase": snakeToLowerCamelCase,
	}).Parse(string(tmpl)))
	if err := t.Execute(out, model); err != nil {
		return err
	}
	return nil
}

var (
	rspace = regexp.MustCompile(`\s+`)
	rsnake = regexp.MustCompile(`_`)
)

func notLast(i, len int) bool {
	return i != len-1
}

func spaceToUpperCamelCase(s string) string {
	if s == "" {
		return ""
	}
	buf := bytes.Buffer{}
	for _, p := range rspace.Split(s, -1) {
		buf.WriteString(strings.ToUpper(p[:1]))
		buf.WriteString(p[1:])
	}
	return buf.String()
}

func snakeToLowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	buf := bytes.Buffer{}
	for i, p := range rsnake.Split(s, -1) {
		if i == 0 {
			buf.WriteString(p)
			continue
		}
		buf.WriteString(strings.ToUpper(p[:1]))
		buf.WriteString(p[1:])
	}
	return buf.String()
}