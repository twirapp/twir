package chat_messages

import (
	"strings"
	"unicode"
)

// TextFuzzyFilter matches message texts against a phrase using the chat wall
// fuzzy semantics: exact case-insensitive substring, or a token whose
// Levenshtein distance (whole token or token prefix of phrase length) to the
// phrase is within MaxDistance. Semantics and thresholds must stay in sync
// with the live matcher in apps/bots/internal/chatwallmatcher.
type TextFuzzyFilter struct {
	Phrase      string
	Length      int
	MaxDistance int
}

const (
	minFuzzyPhraseRunes = 5
	shortPhraseMaxRunes = 8
)

// NewTextFuzzyFilter builds a filter for single-word phrases long enough for
// fuzzy matching and returns nil otherwise (multi-word or short phrases).
func NewTextFuzzyFilter(phrase string) *TextFuzzyFilter {
	tokens := strings.FieldsFunc(
		strings.ToLower(strings.TrimSpace(phrase)),
		func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) },
	)
	if len(tokens) != 1 {
		return nil
	}

	runes := len([]rune(tokens[0]))
	if runes < minFuzzyPhraseRunes {
		return nil
	}

	threshold := 1
	if runes > shortPhraseMaxRunes {
		threshold = 2
	}

	return &TextFuzzyFilter{
		Phrase:      tokens[0],
		Length:      runes,
		MaxDistance: threshold,
	}
}
