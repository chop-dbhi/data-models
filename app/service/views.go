package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func viewIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderIndexHTML(w)
}

func viewModels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	q := r.URL.Query()

	data := map[string]interface{}{
		"Title": "Models",
		"Items": dataModels.List(),
	}

	switch q.Get("format") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelsMarkdown(w, data)
	case "", "html":
		w.Header().Set("content-type", "text/html")
		RenderModelsHTML(w, data)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")

	m := dataModels.Versions(n)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	q := r.URL.Query()

	data := map[string]interface{}{
		"Title": m[0].Name,
		"Items": m,
	}

	switch q.Get("format") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelMarkdown(w, data)
	case "", "html":
		w.Header().Set("content-type", "text/html")
		RenderModelHTML(w, data)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModelVersion(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	m := dataModels.Get(n, v)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	q := r.URL.Query()

	switch q.Get("format") {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelVersionMarkdown(w, m)
	case "", "html":
		w.Header().Set("content-type", "text/html")
		RenderModelVersionHTML(w, m)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
