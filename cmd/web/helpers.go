package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Validation:  make(map[string]string),
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func isValidUrl(input string) bool {
	u, err := url.ParseRequestURI(input)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (app *application) renderComponent(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, err := template.ParseFiles(page)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
