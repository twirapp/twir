// Package chatwallmatcher implements fuzzy phrase matching for chat walls.
//
// Matching semantics are intentionally identical between the live path
// (apps/bots message handler) and the retro sweep SQL used against ClickHouse
// (libs/repositories/chat_messages TextFuzzy filter):
//
//  1. Fast path: case-insensitive substring contains.
//  2. Fuzzy path: per-token Levenshtein distance against the phrase, computed
//     either on the whole token or on the token prefix of phrase length, with
//     a threshold scaled by phrase length.
//
// Single-character substitutions (homoglyphs like a Latin "k" in a Cyrillic
// word, or "0" instead of "о") are caught by the distance itself, so no
// confusable-normalization tables are required anywhere.
package chatwallmatcher

import (
	"strings"
	"unicode"
)

const (
	// minFuzzyPhraseRunes is the minimum phrase length (in runes, without
	// spaces) for fuzzy matching. Shorter phrases only use exact substring
	// matching, because a 1-edit distance on tiny words produces too many
	// false positives.
	minFuzzyPhraseRunes = 5
	// shortPhraseMaxRunes bounds phrases that get a 1-edit threshold.
	// Longer phrases get 2 edits.
	shortPhraseMaxRunes = 8
)

// FuzzyParams describes how a phrase participates in fuzzy matching.
// It is used to parameterize the ClickHouse retro-sweep query with values
// that are computed exactly the same way as for the live matcher.
type FuzzyParams struct {
	// Phrase is the normalized (trimmed, lowercased) phrase. For multi-word
	// phrases the tokens are joined with a single space.
	Phrase string
	// Runes is the rune count of the phrase without spaces. For a single-word
	// phrase it equals the rune count of the only token.
	Runes int
	// Threshold is the maximum allowed Levenshtein distance.
	Threshold int
	// Enabled reports whether fuzzy matching is allowed for this phrase.
	Enabled bool
	// SingleWord reports whether the phrase consists of exactly one token.
	SingleWord bool
}

// Fuzzy computes the fuzzy-matching parameters for a phrase.
func Fuzzy(phrase string) FuzzyParams {
	tokens := splitTokens(strings.ToLower(strings.TrimSpace(phrase)))
	if len(tokens) == 0 {
		return FuzzyParams{}
	}

	total := 0
	for _, token := range tokens {
		total += len([]rune(token))
	}

	return FuzzyParams{
		Phrase:     strings.Join(tokens, " "),
		Runes:      total,
		Threshold:  thresholdForLength(total),
		Enabled:    total >= minFuzzyPhraseRunes,
		SingleWord: len(tokens) == 1,
	}
}

func thresholdForLength(runes int) int {
	if runes > shortPhraseMaxRunes {
		return 2
	}

	return 1
}

// Matcher checks chat text against a chat wall phrase.
type Matcher struct {
	params       FuzzyParams
	phraseTokens []string
}

// New creates a Matcher for the given wall phrase.
func New(phrase string) *Matcher {
	params := Fuzzy(phrase)

	return &Matcher{
		params:       params,
		phraseTokens: splitTokens(params.Phrase),
	}
}

// Matches reports whether the text should be handled by the chat wall.
func (m *Matcher) Matches(text string) bool {
	if m.params.Phrase == "" {
		return false
	}

	lowered := strings.ToLower(text)

	// Fast path: exact substring match, covers natural word extensions
	// ("школьники" for phrase "школьн") and preserves historical behavior.
	if strings.Contains(lowered, m.params.Phrase) {
		return true
	}

	if !m.params.Enabled {
		return false
	}

	tokens := splitTokens(lowered)

	if m.params.SingleWord {
		return m.matchesSingleWord(tokens)
	}

	return m.matchesMultiWord(tokens)
}

func (m *Matcher) matchesSingleWord(tokens []string) bool {
	phrase := m.params.Phrase
	phraseRunes := len([]rune(phrase))

	for _, token := range tokens {
		if levenshtein(token, phrase) <= m.params.Threshold {
			return true
		}

		// Prefix rule: catches evasions hidden inside natural word endings
		// ("шkольницы" -> prefix "шkольн" is 1 edit from "школьн") without
		// matching frames embedded in unrelated words ("прикольных").
		tokenRunes := []rune(token)
		if len(tokenRunes) > phraseRunes {
			if levenshtein(string(tokenRunes[:phraseRunes]), phrase) <= m.params.Threshold {
				return true
			}
		}
	}

	return false
}

func (m *Matcher) matchesMultiWord(tokens []string) bool {
	n := len(m.phraseTokens)
	if len(tokens) < n {
		return false
	}

	for i := 0; i+n <= len(tokens); i++ {
		window := strings.Join(tokens[i:i+n], " ")
		if levenshtein(window, m.params.Phrase) <= m.params.Threshold {
			return true
		}
	}

	return false
}

// splitTokens extracts runs of letters and digits, mirroring the
// extractAll(lowerUTF8(text), '[\p{L}0-9]+') tokenization in ClickHouse.
func splitTokens(s string) []string {
	return strings.FieldsFunc(
		s, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		},
	)
}

// levenshtein computes the rune-based Levenshtein distance between two
// strings, matching ClickHouse editDistanceUTF8 semantics.
func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)

	if len(ar) == 0 {
		return len(br)
	}
	if len(br) == 0 {
		return len(ar)
	}

	prev := make([]int, len(br)+1)
	curr := make([]int, len(br)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(ar); i++ {
		curr[0] = i
		for j := 1; j <= len(br); j++ {
			cost := 0
			if ar[i-1] != br[j-1] {
				cost = 1
			}

			deletion := prev[j] + 1
			insertion := curr[j-1] + 1
			substitution := prev[j-1] + cost

			curr[j] = min(deletion, min(insertion, substitution))
		}
		prev, curr = curr, prev
	}

	return prev[len(br)]
}
