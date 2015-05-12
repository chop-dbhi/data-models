package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/julienschmidt/httprouter"
)

func jsonResponse(w http.ResponseWriter, d interface{}) {
	e := json.NewEncoder(w)

	if err := e.Encode(d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiDataModels(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	models := dataModels.List()

	// In order of model and version.
	sort.Sort(models)

	jsonResponse(w, models)
}

func apiDataModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	// May be bare or have an extension.
	version := p.ByName("version")

	isMarkdown := filepath.Ext(version) == ".md"

	if isMarkdown {
		version = version[:len(version)-3]
	}

	m := dataModels.Get(name, version)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if isMarkdown {
		w.Header().Set("content-type", "text/markdown")
		RenderHierarchyMarkdown(m, w)
	} else {
		jsonResponse(w, m)
	}
}

func apiDataModelSchema(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	version := p.ByName("version")
	ext := p.ByName("ext")

	isMarkdown := ext == ".md"

	m := dataModels.Get(name, version)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if isMarkdown {
		w.Header().Set("content-type", "text/markdown")
		RenderSchemaMarkdown(m, w)
	} else {
		jsonResponse(w, m)
	}
}

func apiDataModelMapping(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	version := p.ByName("version")
	ext := p.ByName("ext")

	isMarkdown := ext == ".md"

	m := dataModels.Get(name, version)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if isMarkdown {
		w.Header().Set("content-type", "text/markdown")
		RenderMappingMarkdown(m, w)
	} else {
		jsonResponse(w, m)
	}
}
func githubWebhook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
