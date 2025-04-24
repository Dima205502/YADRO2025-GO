package words

import (
	"log/slog"
	"maps"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

func Norm(phrase string, logger *slog.Logger) []string {
	words := strings.FieldsFunc(phrase, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	logger.Info("words", "data", words)

	wasStem := make(map[string]bool)
	for _, word := range words {
		if english.IsStopWord(strings.ToLower(word)) {
			continue
		}

		stem := english.Stem(word, false)

		if !wasStem[stem] {
			wasStem[stem] = true
		}
	}

	return slices.Collect(maps.Keys(wasStem))
}
