package main

import (
	"net/http"
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

	// code, header, body := ts.get(t, "/abc123")
	// loc := header.Get("Location")

	// assert.Equal(t, loc, "https://google.com")
	// assert.Equal(t, code, http.StatusSeeOther)
	// assert.StringContains(t, body, `<a href="https://google.com">See Other</a>`)
}
