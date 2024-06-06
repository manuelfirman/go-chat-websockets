package client

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// TemplateHandler struct to represent a template
type TemplateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{
		filename: "chat.html",
	}
}

// ServeHTTP serves the template to the client
func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("pkg/templates", t.filename)))
	})
	t.templ.Execute(w, r)
}
