package main

import (
	"net/http"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "form.tmpl", data)
}

func (app *application) LinkRedirect(w http.ResponseWriter, r *http.Request) {
	suffix := r.URL.Path[1:]
	link, err := app.link.GetBySuffix(r.Context(), suffix)
	if err != nil {

		data := app.newTemplateData(r)
		data.Validation["suffix"] = "Your link is not valid."
		data.Link = &models.Link{}
		app.render(w, r, http.StatusBadRequest, "form.tmpl", data)
		return
	}
	http.Redirect(w, r, link.RedirectUrl, http.StatusSeeOther)
}

func (app *application) LinkPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// TODO: Create a map to store Error messages in to notify user
		app.clientError(w, r, http.StatusBadRequest)
	}
	linkStr := r.PostForm.Get("link")

	data := app.newTemplateData(r)

	// Check if the link is a valid URL
	ok := isValidUrl(linkStr)
	// If the URL is not valid, add an error message to the data map and render the form again
	if !ok {
		data.Validation["url"] = "Invalid Url: Be sure to include 'https://' or 'http://'"
		data.Link = &models.Link{RedirectUrl: linkStr}

		app.render(w, r, http.StatusOK, "form.tmpl", data)
		return
	}

	// Create a new Link struct and generate a suffix
	link := &models.Link{RedirectUrl: linkStr}
	hostAddr := r.Host
	link.Suffix = models.CreateSuffix()

	// Create a short URL
	shortUrl, err := link.CreateShortUrl(hostAddr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	link.ShortUrl = shortUrl

	// Insert the link into the database
	_, err = app.link.Insert(r.Context(), link.RedirectUrl, link.Suffix)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Add the short URL to the data template
	data.Link = link
	// Render the form template with the short URL
	app.render(w, r, http.StatusAccepted, "form.tmpl", data)
}
