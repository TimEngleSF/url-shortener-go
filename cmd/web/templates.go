package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
	"github.com/TimEngleSF/url-shortener-go/ui"
)

type templateData struct {
	CurrentYear     int
	Link            *models.Link
	Validation      map[string]string
	QRImgPath       string
	Form            any
	IsAuthenticated bool
	Flash           string
	ErrorMsg        string
}

/* TEMPLATE FUNCTIONS */
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("Jan 02, 2006 at 3:04PM")
}

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
			"html/partials/*.tmpl",
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
