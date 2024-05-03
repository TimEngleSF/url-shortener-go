package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/assert"
	"github.com/go-playground/form/v4"
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

type mockDestination struct {
	Field1 string `form:"field1"`
	Field2 int    `form:"field2"`
}

func TestDecodePostForm(t *testing.T) {
	app := newTestApplication(t)

	r := httptest.NewRequest("POST", "/test", nil)
	r.PostForm = url.Values{
		"field1": {"value1"},
		"field2": {"123"},
	}

	var dst mockDestination

	err := app.decodePostForm(r, &dst)
	if err != nil {
		t.Errorf("decodePostForm returned an error %v", err)
	}

	if dst.Field1 != "value1" {
		t.Errorf("Expected dst.Field1 to be 'value1', got %s", dst.Field1)
	}
	if dst.Field2 != 123 {
		t.Errorf("Expected dst.Field2 to be 123, got %d", dst.Field2)
	}
}

func TestDecodePostFormErrors(t *testing.T) {
	app := newTestApplication(t)

	// Test case 1: Simulate  decoder error
	r1 := httptest.NewRequest("POST", "/test", nil)
	r1.PostForm = url.Values{
		"field1": {"value1"},
		"field2": {"123"},
	}
	// Mock form decoder to return an error
	app.formDecoder = &mockDecoder{err: errors.New("form decoder error")}
	var dst1 mockDestination
	err1 := app.decodePostForm(r1, &dst1)
	if err1 == nil || err1.Error() != "form decoder error" {
		t.Errorf("decodePostForm did not return expected error")
	}

	// Test case 2: Simulate HTTP request error
	r2 := httptest.NewRequest("POST", "/test", nil)
	// Set an invalid form data format to simulate an error
	r2.PostForm = url.Values{"invalid": {"data"}}
	var dst2 mockDestination
	err2 := app.decodePostForm(r2, &dst2)
	if err2 == nil {
		t.Errorf("decodePostForm did not return an error for invalid form data")
	}

	// Test case 3: Simulate panic due to invalid decoder error
	r3 := httptest.NewRequest("POST", "/test", nil)
	r3.PostForm = url.Values{"field1": {"value1"}}
	// Mock form decoder to return an invalid decoder error
	app.formDecoder = &mockDecoder{err: &form.InvalidDecoderError{}}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("decodePostForm did not panic as expected for invalid decoder error")
		}
	}()
	var dst3 mockDestination
	_ = app.decodePostForm(r3, &dst3)
}

type mockDecoder struct {
	err error
}

func (d *mockDecoder) Decode(dst interface{}, values url.Values) error {
	return d.err
}
