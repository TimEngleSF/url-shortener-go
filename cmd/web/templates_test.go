package main

import (
	"testing"
	"time"

	"github.com/TimEngleSF/url-shortener-go/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC AM",
			tm:   time.Date(2024, 4, 22, 11, 50, 0, 0, time.UTC),
			want: "Apr 22, 2024 at 11:50AM",
		},
		{
			name: "UTC PM",
			tm:   time.Date(2024, 4, 22, 15, 50, 0, 0, time.UTC),
			want: "Apr 22, 2024 at 3:50PM",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 4, 22, 15, 50, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "Apr 22, 2024 at 2:50PM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
