package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func viewModels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch p.ByName("ext") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelsMarkdown(w, dataModels)
		break
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderModelsHTML(w, dataModels)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModelFull(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	m := dataModels.Get(n, v)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch p.ByName("ext") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderFullMarkdown(w, m)
		break
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderFullHTML(w, m)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModelDefinition(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	m := dataModels.Get(n, v)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch p.ByName("ext") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderDefinitionMarkdown(w, m)
		break
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderSchemaHTML(w, m)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModelSchema(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	m := dataModels.Get(n, v)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch p.ByName("ext") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderSchemaMarkdown(w, m)
		break
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderSchemaHTML(w, m)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModelMapping(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	m := dataModels.Get(n, v)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch p.ByName("ext") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderMappingMarkdown(w, m)
		break
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderMappingHTML(w, m)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
