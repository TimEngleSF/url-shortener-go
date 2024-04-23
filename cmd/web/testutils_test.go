package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache: templateCache,
		link:          &mocks.LinkMock{},
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
