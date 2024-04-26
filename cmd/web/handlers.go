package main

import (
	"errors"
	"net/http"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "form.tmpl", data)
}

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// ////////////////////// LINKS ////////////////////////
//
// // REDIRECT ////
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

// // LINK POST ////
func (app *application) LinkPost(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	var err error
	data := app.newTemplateData(r)

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
	}

	linkStr := r.PostForm.Get("link")

	// Check if the link is a valid URL
	ok := isValidUrl(linkStr)
	// If the URL is not valid, add an error message to the data map and render the form again
	if !ok {
		data.Validation["url"] = "Invalid Url: Be sure to include 'https://' or 'http://'"
		data.Link = &models.Link{RedirectUrl: linkStr}

		app.render(w, r, http.StatusOK, "form.tmpl", data)
		return
	}

	// Check if the URL already exists in the database
	link, err = app.link.GetByURL(r.Context(), linkStr)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// If the URL does not exist, create a new Link struct
			// in the next block of code
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	// If the URL does not exist in the database, create a new Link struct
	if link.Suffix == "" {
		link = models.Link{RedirectUrl: linkStr}
		link.Suffix = models.CreateSuffix()

		// Insert the link into the database
		_, err = app.link.Insert(r.Context(), link.RedirectUrl, link.Suffix)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	// Create a short URL
	hostAddr := r.Host
	shortUrl, err := link.CreateShortUrl(hostAddr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	link.ShortUrl = shortUrl
	// Add the short URL to the data template
	data.Link = &link

	// Create QR
	qrPath, err := app.qr.CreateMedium(shortUrl)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.QRImgPath = qrPath

	// Render the form template with the short URL
	app.render(w, r, http.StatusCreated, "form.tmpl", data)
}

// ////////////////////// USERS ////////////////////////
//
// // SIGNUP FORM ////
func (app *application) SignUpForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signup Form"))
}

// // LOGIN FORM ////
func (app *application) LoginForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login Form"))
}
