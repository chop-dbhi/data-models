package main

import (
	"io"
	"text/template"
)

var modelTemplate *template.Template

func init() {
	data, err := Asset("templates/model.md")

	if err != nil {
		panic(err)
	}

	modelTemplate = template.Must(template.New("model").Parse(string(data)))
}

func RenderMarkdownHierarchy(m *Model, w io.Writer) {
	modelTemplate.Execute(w, m)
}
