package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/assert"
)

func TestServerError(t *testing.T) {
	app := newTestApplication(t)
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET req",
			method: "GET",
			path:   "/path",
		},
		{
			name:   "POST req",
			method: "POST",
			path:   "/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)

			rr := httptest.NewRecorder()

			testErr := errors.New("Test Error")
			app.serverError(rr, req, testErr)

			trimmedBody := strings.TrimSpace(rr.Body.String())
			expectedBody := http.StatusText(http.StatusInternalServerError)

			assert.Equal(t, rr.Code, http.StatusInternalServerError)
			assert.StringsContains(t, expectedBody, []string{trimmedBody})
		})
	}
}

func TestClientError(t *testing.T) {
	tests := []struct {
		name   string
		method string
		code   int
		want   int
	}{
		{
			name:   "POST Unprocessable Entity",
			method: "POST",
			code:   http.StatusUnprocessableEntity,
			want:   422,
		},
		{
			name:   "GET Bad Request",
			method: "GET",
			code:   http.StatusBadRequest,
			want:   400,
		},
	}

	app := newTestApplication(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/invalid", nil)
			rr := httptest.NewRecorder()

			app.clientError(rr, req, tt.code)

			assert.Equal(t, rr.Code, tt.want)
		})
	}
}
