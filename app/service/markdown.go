package main

import (
	"io"
	"text/template"
)

func loadTemplate(n string) *template.Template {
	data, err := Asset(n)

	if err != nil {
		panic(err)
	}

	return template.Must(template.New("model").Parse(string(data)))
}

func RenderModelsMarkdown(w io.Writer, v interface{}) {
	t := loadTemplate("templates/models.md")
	t.Execute(w, v)
}

func RenderFullMarkdown(w io.Writer, v interface{}) {
	t := loadTemplate("templates/full.md")
	t.Execute(w, v)
}

func RenderDefinitionMarkdown(w io.Writer, v interface{}) {
	t := loadTemplate("templates/definition.md")
	t.Execute(w, v)
}

func RenderSchemaMarkdown(w io.Writer, v interface{}) {
	t := loadTemplate("templates/schema.md")
	t.Execute(w, v)
}

func RenderMappingMarkdown(w io.Writer, v interface{}) {
	t := loadTemplate("templates/mappings.md")
	t.Execute(w, v)
}
