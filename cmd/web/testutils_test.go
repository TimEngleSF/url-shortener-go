package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/models/mocks"
	"github.com/TimEngleSF/url-shortener-go/internal/qr"
	"github.com/go-playground/form/v4"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	return &application{
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache: templateCache,
		link:          &mocks.LinkMock{},
		qr:            &qr.QRCodeMock{},
		formDecoder:   formDecoder,
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
