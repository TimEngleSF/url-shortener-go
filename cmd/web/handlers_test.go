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

	code, header, body := ts.get(t, "/abc123")
	loc := header.Get("Location")

	assert.Equal(t, loc, "https://google.com")
	assert.Equal(t, code, http.StatusSeeOther)
	assert.StringContains(t, body, `<a href="https://google.com">See Other</a>`)
}
