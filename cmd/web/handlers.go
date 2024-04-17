package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
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
		data.Link = &models.Link{RedirectUrl: linkStr}

		err = ts.ExecuteTemplate(w, "form", data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	link := &models.Link{RedirectUrl: linkStr}
	hostAddr := r.Host
	link.Suffix = models.CreateSuffix()

	shortUrl, err := link.CreateShortUrl(hostAddr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	link.ShortUrl = "https://" + shortUrl
	_, err = app.link.Insert(r.Context(), link.RedirectUrl, link.Suffix)
	if err != nil {
		fmt.Println("Error inserting link into database: ", err)
		app.serverError(w, r, err)
		return
	}
	data.Link = link

	err = ts.ExecuteTemplate(w, "form", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
