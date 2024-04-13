package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/form/form.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) LinkPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// TODO: Create a map to store Error messages in to notify user
		app.clientError(w, r, http.StatusBadRequest)
	}
	linkStr := r.PostForm.Get("link")

	files := []string{
		"./ui/html/form/form.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)

	ok := isValidUrl(linkStr)
	if !ok {
		data.Validation["url"] = "Invalid Url: Be sure to include 'https://' or 'http://'"
		data.Link = &Link{RedirectUrl: linkStr}

		err = ts.ExecuteTemplate(w, "form", data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	link := &Link{RedirectUrl: linkStr}
	data.Link = link
	err = ts.ExecuteTemplate(w, "form", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
