package cli

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "GO POKE",
			expected: []string{"go", "poke"},
		},
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d", c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q) == %q, expected %q", c.input, word, expectedWord)
			}
		}
	}
}

