package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/assert"
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
			wantCode:     400,
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
	ts := newTestServer(t, app.routes())
	host := ts.URL
	host, _ = strings.CutPrefix(host, "http")

	tests := []struct {
		name         string
		redirectURL  string
		wantSuffix   string
		wantShortUrl string
		wantCode     int
		displayText  []string
	}{
		{
			name:        "Valid Existing URL",
			redirectURL: "https://google.com",
			displayText: []string{
				fmt.Sprintf("%s/abc123", host),
			},
			wantCode: http.StatusCreated,
		},
		{
			name:        "Invalid URL",
			redirectURL: "google",
			wantSuffix:  "",
			displayText: []string{
				"Invalid Url: Be sure to include",
				"google",
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "Valid New URL",
			redirectURL: "https://yahoo.com",
			displayText: []string{host},
			wantCode:    http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("link", tt.redirectURL)
			code, _, body := ts.postForm(t, "/link/new", form)

			assert.StringsContains(t, body, tt.displayText)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}
