package main

import (
	"io"
	"text/template"
)

var (
	hierarchyTemplate, schemaTemplate, mappingsTemplate *template.Template
)

func loadTemplate(n string) *template.Template {
	data, err := Asset(n)

	if err != nil {
		panic(err)
	}

	return template.Must(template.New("model").Parse(string(data)))
}

func init() {
	hierarchyTemplate = loadTemplate("templates/model.md")
	schemaTemplate = loadTemplate("templates/schema.md")
	mappingsTemplate = loadTemplate("templates/mappings.md")
}

func RenderHierarchyMarkdown(m *Model, w io.Writer) {
	hierarchyTemplate.Execute(w, m)
}

func RenderSchemaMarkdown(m *Model, w io.Writer) {
	schemaTemplate.Execute(w, m)
}

func RenderMappingMarkdown(m *Model, w io.Writer) {
	mappingsTemplate.Execute(w, m)
}
