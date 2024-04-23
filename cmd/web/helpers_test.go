package main

import (
	"testing"

	"github.com/TimEngleSF/url-shortener-go/internal/assert"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid http",
			url:  "http://google.com",
			want: true,
		},
		{
			name: "Invalid http",
			url:  "http:google.com",
			want: false,
		},
		{
			name: "Valid https",
			url:  "https://google.com",
			want: true,
		},
		{
			name: "Invalid https",
			url:  "https/google.com",
			want: false,
		},
		{
			name: "Missing http or https",
			url:  "google.com",
			want: false,
		},
		{
			name: "direct path",
			url:  "/foo/bar",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := isValidUrl(tt.url)

			assert.Equal(t, isValid, tt.want)
		})
	}
}
