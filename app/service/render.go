package main

import (
	"bytes"
	"io"
	"text/template"

	"github.com/russross/blackfriday"
)

const (
	htmlFlags = blackfriday.HTML_SKIP_HTML |
		blackfriday.HTML_USE_SMARTYPANTS

	extFlags = blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_FENCED_CODE
)

func loadTemplate(n string) *template.Template {
	data, err := Asset(n)

	if err != nil {
		panic(err)
	}

	return template.Must(template.New("model").Parse(string(data)))
}

func renderMarkdown(w io.Writer, n string, v interface{}) {
	t := loadTemplate(n)
	t.Execute(w, v)
}

func renderHTML(w io.Writer, b []byte) {
	// Render the final HTML page.
	t := loadTemplate("assets/wrap.html")
	s, _ := Asset("assets/style.css")

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	c := blackfriday.Markdown(b, renderer, extFlags)

	data := map[string]string{
		"Content": string(c),
		"Style":   string(s),
	}

	t.Execute(w, data)
}

func RenderModelsMarkdown(w io.Writer, v interface{}) {
	renderMarkdown(w, "assets/models.md", v)
}

func RenderModelMarkdown(w io.Writer, v interface{}) {
	renderMarkdown(w, "assets/models.md", v)
}

func RenderModelVersionMarkdown(w io.Writer, v interface{}) {
	renderMarkdown(w, "assets/full.md", v)
}

func RenderModelCompareMarkdown(w io.Writer, m1 *Model, m2 *Model) {
	DiffModels(w, m1, m2)
}

func RenderReposMarkdown(w io.Writer, v interface{}) {
	renderMarkdown(w, "assets/repos.md", v)
}

func RenderIndexHTML(w io.Writer) {
	b := bytes.Buffer{}
	renderMarkdown(&b, "assets/index.md", nil)
	renderHTML(w, b.Bytes())
}

func RenderModelsHTML(w io.Writer, v interface{}) {
	b := bytes.Buffer{}
	renderMarkdown(&b, "assets/models.md", v)
	renderHTML(w, b.Bytes())
}

func RenderModelHTML(w io.Writer, v interface{}) {
	b := bytes.Buffer{}
	renderMarkdown(&b, "assets/models.md", v)
	renderHTML(w, b.Bytes())
}

func RenderModelVersionHTML(w io.Writer, v interface{}) {
	b := bytes.Buffer{}
	renderMarkdown(&b, "assets/full.md", v)
	renderHTML(w, b.Bytes())
}

func RenderModelCompareHTML(w io.Writer, m1 *Model, m2 *Model) {
	b := bytes.Buffer{}
	DiffModels(&b, m1, m2)
	renderHTML(w, b.Bytes())
}

func RenderReposHTML(w io.Writer, v interface{}) {
	b := bytes.Buffer{}
	RenderReposMarkdown(&b, v)
	renderHTML(w, b.Bytes())
}
