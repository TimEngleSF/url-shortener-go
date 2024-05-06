package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/TimEngleSF/url-shortener-go/internal/models/mocks"
	"github.com/TimEngleSF/url-shortener-go/internal/qr"
	S3 "github.com/TimEngleSF/url-shortener-go/internal/s3"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache:  templateCache,
		link:           &mocks.LinkMock{},
		user:           &mocks.UserMock{},
		qr:             &qr.QRCodeMock{},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		s3:             &S3.S3Mock{},
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (status int, header http.Header, body string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = string(bytes.TrimSpace(b))

	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (code int, header http.Header, body string) {
	path := ts.URL + urlPath
	rs, err := ts.Client().PostForm(path, form)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = string(bytes.TrimSpace(b))

	return rs.StatusCode, rs.Header, body
}

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'/>`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}
