package words

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNorm(t *testing.T) {
	testCases := []struct {
		name     string
		given    string
		expected []string
	}{
		{
			name:     "empty",
			given:    "",
			expected: []string{},
		},
		{
			name:     "simple",
			given:    "simple",
			expected: []string{"simpl"},
		},
		{
			name:     "followers",
			given:    "I follow followers",
			expected: []string{"follow"},
		},
		{
			name:     "punctuation",
			given:    "I shouted: 'give me your car!!!",
			expected: []string{"shout", "give", "car"},
		},
		{
			name:     "stop words only",
			given:    "I and you or me or them, who will?",
			expected: []string{},
		},
		{
			name:     "weird",
			given:    "Moscow!123'check-it'or   123, man,that,difficult:heck",
			expected: []string{"moscow", "check", "123", "man", "difficult", "heck"},
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := Norm(tc.given, logger)
			assert.ElementsMatch(t, tc.expected, actual)
		})
	}

}
