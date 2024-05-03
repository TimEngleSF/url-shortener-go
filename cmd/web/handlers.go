package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
	validator "github.com/TimEngleSF/url-shortener-go/internal/validators"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = linkCreateForm{}

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
		data.Form = linkCreateForm{}
		data.Validation["suffix"] = "Your link is not valid."
		data.Link = &models.Link{}
		app.render(w, r, http.StatusBadRequest, "form.tmpl", data)
		return
	}
	http.Redirect(w, r, link.RedirectUrl, http.StatusSeeOther)
}

// // LINK POST ////

type linkCreateForm struct {
	Link models.Link
	validator.Validator
}

func (app *application) LinkPost(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	var err error
	data := app.newTemplateData(r)

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
	}

	linkStr := r.PostForm.Get("link")

	form := linkCreateForm{
		Link: models.Link{
			RedirectUrl: linkStr,
		},
	}

	form.CheckField(validator.NotBlank(form.Link.RedirectUrl), "url", "Invalid URL: URL must not be blank")
	form.CheckField(validator.IsValidUrl(form.Link.RedirectUrl), "url", "Invalid URL: Be sure to include 'https://' or 'http://'")
	form.CheckField(validator.MaxChars(form.Link.RedirectUrl, 500), "url", "Invalid URL: URL length too long")

	data.Form = form

	if !form.Valid() {
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
	data.Flash = "Link successfully created!"
	app.render(w, r, http.StatusCreated, "form.tmpl", data)
}

// ////////////////////// USERS ////////////////////////
//
// // SIGNUP FORM ////
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userAddForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

type userAddForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userAddForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Name, 3), "name", "Name must be at least 3 characters long")
	form.CheckField(validator.MaxChars(form.Name, 20), "name", "Name must be 20 characters or less")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsValidEmail(form.Email), "email", "This is not a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")

	// Init data & set form values
	data := app.newTemplateData(r)
	data.Form = form

	// Return and display form field errors
	if !form.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Return and display err if email in use
	exists, _ := app.user.Exists(r.Context(), form.Email)
	if exists {
		form.AddFieldError("email", "This Email already in use")
		data.Form = form
		return
	}

	// Insert account into db
	err = app.user.Insert(r.Context(), form.Name, strings.ToLower(form.Email), form.Password)
	// Return and display db query errors
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email already in use")
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			data.ErrorMsg = "Error creating account"
			app.render(w, r, http.StatusInternalServerError, "signup.tmpl", data)
		}
		return
	}

	// Account created successfully
	app.sessionManager.Put(r.Context(), "flash", "Successfully created account")
	data.Form = userAddForm{}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// // LOGIN FORM ////
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		data.ErrorMsg = "There was an error logging in"
		app.render(w, r, http.StatusInternalServerError, "login.tmpl", data)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsValidEmail(form.Email), "email", "This is not a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")

	data.Form = form

	if !form.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Authenticate user
	user, err := app.user.Authenticate(r.Context(), form.Email, form.Password)
	// Return and display error if credentials are invalid or err querying db
	if err != nil {
		var status int
		if errors.Is(err, models.ErrInvalidCredentials) {
			data.ErrorMsg = "Email or password invalid"
			status = http.StatusUnauthorized
		} else {
			data.ErrorMsg = "There was an error accessing account"
			status = http.StatusInternalServerError
		}
		app.render(w, r, status, "login.tmpl", data)
		return
	}

	// TODO:
	// If authenticated need to authorize user & update sessionManager
	fmt.Fprint(w, user.Email)
}

//// LOGOUT FORM ////

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logout user")
}
