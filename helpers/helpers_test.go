package helpers

import (
	"testing"
)

func TestGenerateParam(t *testing.T) {
	subtests := []int{0, 1, 2, 3, 4, 5, 6, 7}

	for _, st := range subtests {
		actual := GenerateParam(st)

		if len(actual) != st {
			t.Errorf("Output mistmatch. Expected %s to have length of %d, but got %d\n", actual, st, len(actual))
		}
	}
}

func TestAddURLPrefix( t *testing.T) {
	subtests := []struct {
		expected string
		input    string
	}{
		{
			expected: "https://google.com",
			input:    "google.com",
		},
		{
			expected: "https://www.google.com",
			input:    "www.google.com",
		},
		{
			expected: "http://google.com",
			input: "http://google.com",
		},
		{
			expected: "https://google.com",
			input: "https://google.com",
		},
	}

	for _, st := range subtests {
		actual := AddURLPrefix(st.input)

		if actual != st.expected {
			t.Errorf("Output mistmatch. Expected %s but got %s\n", st.expected, actual)
		}
	}
}
