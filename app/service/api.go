package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func jsonResponse(w http.ResponseWriter, d interface{}) {
	e := json.NewEncoder(w)

	if err := e.Encode(d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiModels(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	models := dataModels.List()
	jsonResponse(w, models)
}

func apiModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")

	var m Models

	if m = dataModels.Versions(n); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse(w, m)
}

func apiModelVersion(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	var m *Model

	if m = dataModels.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse(w, m)
}

func apiTable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")
	tn := p.ByName("table")

	var (
		m *Model
		t *Table
	)

	if m = dataModels.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t = m.Tables.Get(tn); t == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse(w, t)
}

func apiField(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")
	tn := p.ByName("table")
	fn := p.ByName("field")

	var (
		m *Model
		t *Table
		f *Field
	)

	if m = dataModels.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t = m.Tables.Get(tn); t == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if f = t.Fields.Get(fn); f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse(w, f)
}
