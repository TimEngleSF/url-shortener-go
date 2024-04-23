package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got %v; want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func StringsContains(t *testing.T, actual string, expectedSubstrings []string) {
	t.Helper()

	for _, s := range expectedSubstrings {
		if !strings.Contains(actual, s) {
			t.Errorf("got: %q; expected to contain: %q", actual, s)
		}
	}
}
