package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/assert"
	"github.com/TimEngleSF/url-shortener-go/internal/qr"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/")

	button := `<li><button type="submit">Get Link</button></li>`

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, button)
}

func TestLinkRedirect(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name           string
		suffix         string
		wantLocation   string
		wantCode       int
		wantBody       string
		wantValidation string
		err            *error
	}{
		{
			name:         "Valid redirect",
			suffix:       "/abc123",
			wantLocation: "https://google.com",
			wantCode:     http.StatusSeeOther,
			wantBody:     `<a href="https://google.com">See Other</a>`,
		},
		{
			name:         "Invalid redirect",
			suffix:       "/aaaaaa",
			wantLocation: "",
			wantCode:     http.StatusBadRequest,
			wantBody:     `Your link is not valid.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, header, body := ts.get(t, tt.suffix)
			loc := header.Get("Location")
			assert.Equal(t, code, tt.wantCode)
			assert.Equal(t, loc, tt.wantLocation)
			assert.StringContains(t, body, tt.wantBody)
		})
	}
}

func TestLinkPost(t *testing.T) {
	app := newTestApplication(t)
	app.qr = &qr.QRCodeMock{}

	ts := newTestServer(t, app.routes())
	host := ts.URL
	host, _ = strings.CutPrefix(host, "http")

	_, _, body := ts.get(t, "/")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		redirectURL  string
		wantSuffix   string
		wantShortUrl string
		wantCode     int
		csrfToken    string
		displayText  []string
		qrCalled     bool
	}{
		{
			name:        "Valid Existing URL",
			redirectURL: "https://google.com",
			csrfToken:   validCSRFToken,
			displayText: []string{
				fmt.Sprintf("%s/abc123", host),
			},
			wantCode: http.StatusCreated,
			qrCalled: true,
		},
		{
			name:        "Invalid URL",
			redirectURL: "google",
			csrfToken:   validCSRFToken,
			wantSuffix:  "",
			displayText: []string{
				"Invalid URL: Be sure to include",
				"google",
			},
			wantCode: http.StatusOK,
			qrCalled: false,
		},
		{
			name:        "Valid New URL",
			redirectURL: "https://yahoo.com",
			csrfToken:   validCSRFToken,
			displayText: []string{host},
			wantCode:    http.StatusCreated,
			qrCalled:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("csrtf token %s", tt.csrfToken)
			form := url.Values{}
			form.Add("link", tt.redirectURL)
			form.Add("csrf_token", tt.csrfToken)
			code, _, body := ts.postForm(t, "/link/new", form)

			assert.StringsContains(t, body, tt.displayText)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}
