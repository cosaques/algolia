package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// templateHandler allows to handle the static html-templates.
type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

// ServeHTTP implements http.Handler interface.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	t.templ.Execute(w, r)
}
