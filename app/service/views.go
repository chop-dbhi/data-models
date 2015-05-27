package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	defaultFormat = "html"

	mimetypes = map[string]string{
		"text/markdown":    "markdown",
		"text/html":        "html",
		"application/json": "json",
	}

	queryFormats = map[string]string{
		"md":       "markdown",
		"markdown": "markdown",
		"html":     "html",
		"json":     "json",
	}
)

func jsonResponse(w http.ResponseWriter, d interface{}) {
	e := json.NewEncoder(w)

	if err := e.Encode(d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// detectFormat applies content negotiation logic to determine the
// appropriate response representation.
func detectFormat(w http.ResponseWriter, r *http.Request) string {
	var (
		ok     bool
		format string
	)

	format = queryFormats[strings.ToLower(r.URL.Query().Get("format"))]

	// Query parameter
	if format == "" {
		// Accept header
		acceptType := r.Header.Get("Accept")
		acceptType, _, _ = mime.ParseMediaType(acceptType)

		// Fallback to default
		if format, ok = mimetypes[acceptType]; !ok {
			format = defaultFormat
		}
	}

	var contentType string

	switch format {
	case "html":
		contentType = "text/html"
	case "markdown":
		contentType = "text/markdown"
	case "json":
		contentType = "application/json"
	}

	w.Header().Set("content-type", contentType)

	return format
}

func viewIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch detectFormat(w, r) {
	case "html":
		RenderIndexHTML(w)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func viewModels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := map[string]interface{}{
		"Title": "Models",
		"Items": dataModels.List(),
	}

	switch detectFormat(w, r) {
	case "markdown":
		RenderModelsMarkdown(w, data)
	case "html":
		RenderModelsHTML(w, data)
	case "json":
		jsonResponse(w, data["Items"])
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func viewModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")

	m := dataModels.Versions(n)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": m[0].Name,
		"Items": m,
	}

	switch detectFormat(w, r) {
	case "markdown":
		RenderModelMarkdown(w, data)
	case "html":
		RenderModelHTML(w, data)
	case "json":
		jsonResponse(w, m)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
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

	switch detectFormat(w, r) {
	case "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelVersionMarkdown(w, m)
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderModelVersionHTML(w, m)
	case "json":
		jsonResponse(w, m)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func viewTable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	switch detectFormat(w, r) {
	case "json":
		jsonResponse(w, t)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func viewField(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	switch detectFormat(w, r) {
	case "json":
		jsonResponse(w, f)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func viewCompareModels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n1 := p.ByName("name1")
	v1 := p.ByName("version1")

	m1 := dataModels.Get(n1, v1)

	if m1 == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	n2 := p.ByName("name2")
	v2 := p.ByName("version2")

	m2 := dataModels.Get(n2, v2)

	if m2 == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch detectFormat(w, r) {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelCompareMarkdown(w, m1, m2)
	case "", "html":
		w.Header().Set("content-type", "text/html")
		RenderModelCompareHTML(w, m1, m2)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func verifyGithubSignature(sig string, r io.Reader) bool {
	mac := hmac.New(sha1.New, []byte(secret))
	io.Copy(mac, r)
	expected := fmt.Sprintf("sha1=%x", mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sig))
}

func viewUpdateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// If no secret has been supplied, update at will.
	if secret == "" {
		updateRepo()
		return
	}

	// Check for Github's webhook signature.
	if sig := r.Header.Get("X-Hub-Signature"); sig != "" {
		defer r.Body.Close()

		if verifyGithubSignature(sig, r.Body) {
			updateRepo()
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func viewModelSchema(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	var (
		m *Model
	)

	if m = dataModels.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	aux := make(map[string]interface{})

	aux["model"] = m.Name
	aux["version"] = m.Version
	aux["tables"] = m.Tables
	aux["schema"] = m.schema

	switch detectFormat(w, r) {
	case "json":
		jsonResponse(w, aux)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
