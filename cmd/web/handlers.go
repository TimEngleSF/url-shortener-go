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

		app.clientError(w, r, http.StatusNotFound)
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

	ok := isValidUrl(linkStr)
	if !ok {
		data.Validation["url"] = "Invalid Url: Be sure to include 'https://' or 'http://'"
		data.Link = &models.Link{RedirectUrl: linkStr}

		app.render(w, r, http.StatusOK, "form.tmpl", data)
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

	link.ShortUrl = shortUrl

	_, err = app.link.Insert(r.Context(), link.RedirectUrl, link.Suffix)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.Link = link

	app.render(w, r, http.StatusAccepted, "form.tmpl", data)
}
