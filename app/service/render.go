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
		blackfriday.EXTENSION_HEADER_IDS
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

func renderHTML(w io.Writer, n string, v interface{}) {
	t := loadTemplate(n)
	b := bytes.Buffer{}

	t.Execute(&b, v)

	// Render the final HTML page.
	t = loadTemplate("assets/wrap.html")
	s, _ := Asset("assets/style.css")

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	c := blackfriday.Markdown(b.Bytes(), renderer, extFlags)

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

func RenderIndexHTML(w io.Writer) {
	renderHTML(w, "assets/index.md", nil)
}

func RenderModelsHTML(w io.Writer, v interface{}) {
	renderHTML(w, "assets/models.md", v)
}

func RenderModelHTML(w io.Writer, v interface{}) {
	renderHTML(w, "assets/models.md", v)
}

func RenderModelVersionHTML(w io.Writer, v interface{}) {
	renderHTML(w, "assets/full.md", v)
}
