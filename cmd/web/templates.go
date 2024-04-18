package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
	"github.com/TimEngleSF/url-shortener-go/ui"
)

type templateData struct {
	CurrentYear int
	Link        *models.Link
	Validation  map[string]string
}

/* TEMPLATE FUNCTIONS */
var functions = template.FuncMap{}

/* TEMPLATE CACHE */
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl",
			// "html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
