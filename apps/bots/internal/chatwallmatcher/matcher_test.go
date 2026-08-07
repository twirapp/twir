package chatwallmatcher

import (
	"testing"
)

func TestMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		phrase string
		text   string
		want   bool
	}{
		// Exact / fast path
		{name: "exact", phrase: "школьн", text: "школьн", want: true},
		{name: "exact case insensitive", phrase: "школьн", text: "ШКОЛЬН", want: true},
		{name: "phrase case insensitive", phrase: "Школьн", text: "школьн", want: true},
		{name: "word extension contains", phrase: "школьн", text: "школьники", want: true},
		{name: "word extension feminine", phrase: "школьн", text: "школьницы", want: true},
		{name: "in a sentence", phrase: "школьн", text: "кроме школьниц есть еще боты", want: true},

		// Evasions caught by fuzzy path
		{name: "digit substitution", phrase: "школьн", text: "шк0льн", want: true},
		{name: "latin k substitution", phrase: "школьн", text: "шkольн", want: true},
		{name: "latin k with extension", phrase: "школьн", text: "шkольницы", want: true},
		{name: "truncated", phrase: "школьн", text: "школь", want: true},
		{name: "extra letter suffix", phrase: "школьн", text: "школьны", want: true},
		{name: "extra letter suffix in sentence", phrase: "школьн", text: "вот школьны", want: true},
		{name: "uppercase extension", phrase: "школьн", text: "Школьницы", want: true},
		{name: "latin substitution mid sentence", phrase: "blocked", text: "you are bl0cked now", want: true},
		{name: "long phrase two edits", phrase: "школьницы", text: "шкoльнiцы", want: true},

		// Lookalikes that must NOT match
		{name: "shkola", phrase: "школьн", text: "школа", want: false},
		{name: "skol", phrase: "школьн", text: "сколь", want: false},
		{name: "shkol", phrase: "школьн", text: "школ", want: false},
		{name: "embedded frame in unrelated word", phrase: "школьн", text: "прикольных", want: false},
		{name: "first letter substitution evasion", phrase: "школьн", text: "хкольн", want: true},
		{name: "first letter substitution extended", phrase: "школьн", text: "хкольницы", want: true},
		// Accepted tradeoff: a 1-edit prefix of a common word also matches.
		{name: "skolko prefix tradeoff", phrase: "школь", text: "а сколько у тебя", want: true},
		{name: "unrelated", phrase: "школьн", text: "привет", want: false},
		{name: "empty text", phrase: "школьн", text: "", want: false},
		{name: "partial tokens nowhere near", phrase: "школьн", text: "ш к о ль н", want: false},

		// Multi-word phrases
		{name: "multiword exact", phrase: "голые школьницы", text: "голые школьницы", want: true},
		{name: "multiword in sentence", phrase: "голые школьницы", text: "вот голые школьницы тут", want: true},
		{name: "multiword fuzzy", phrase: "голые школьницы", text: "голые шк0льницы", want: true},
		{name: "multiword partial only", phrase: "голые школьницы", text: "голые", want: false},
		{name: "multiword reversed", phrase: "голые школьницы", text: "школьницы голые", want: false},

		// Short phrases: fuzzy disabled, exact contains only
		{name: "short phrase exact", phrase: "школ", text: "школы", want: true},
		{name: "short phrase no fuzzy", phrase: "школ", text: "шкoл", want: false},

		// Edge cases
		{name: "empty phrase", phrase: "", text: "школьн", want: false},
		{name: "blank phrase", phrase: "   ", text: "школьн", want: false},
		{name: "phrase with punctuation", phrase: "школьн!", text: "школьн!11", want: true},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				if got := New(tt.phrase).Matches(tt.text); got != tt.want {
					t.Errorf("Matches(%q, %q) = %v, want %v", tt.text, tt.phrase, got, tt.want)
				}
			},
		)
	}
}

func TestFuzzyParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		phrase     string
		wantPhrase string
		wantRunes  int
		wantThr    int
		wantOn     bool
		wantSingle bool
	}{
		{phrase: "школьн", wantPhrase: "школьн", wantRunes: 6, wantThr: 1, wantOn: true, wantSingle: true},
		{phrase: "Школьн", wantPhrase: "школьн", wantRunes: 6, wantThr: 1, wantOn: true, wantSingle: true},
		{phrase: "  школьн  ", wantPhrase: "школьн", wantRunes: 6, wantThr: 1, wantOn: true, wantSingle: true},
		{phrase: "школ", wantPhrase: "школ", wantRunes: 4, wantThr: 1, wantOn: false, wantSingle: true},
		{phrase: "школьницы", wantPhrase: "школьницы", wantRunes: 9, wantThr: 2, wantOn: true, wantSingle: true},
		{phrase: "голые школьницы", wantPhrase: "голые школьницы", wantRunes: 14, wantThr: 2, wantOn: true, wantSingle: false},
		{phrase: "", wantOn: false, wantSingle: false},
	}

	for _, tt := range tests {
		t.Run(
			tt.phrase, func(t *testing.T) {
				t.Parallel()

				got := Fuzzy(tt.phrase)
				if got.Phrase != tt.wantPhrase ||
					got.Runes != tt.wantRunes ||
					got.Threshold != tt.wantThr ||
					got.Enabled != tt.wantOn ||
					got.SingleWord != tt.wantSingle {
					t.Errorf("Fuzzy(%q) = %+v, want phrase=%q runes=%d thr=%d enabled=%v single=%v",
						tt.phrase, got, tt.wantPhrase, tt.wantRunes, tt.wantThr, tt.wantOn, tt.wantSingle)
				}
			},
		)
	}
}

func TestLevenshtein(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a, b string
		want int
	}{
		{"школьн", "школьн", 0},
		{"шк0льн", "школьн", 1},
		{"шkольн", "школьн", 1}, // latin k: 1 rune substitution, not 2 bytes
		{"школь", "школьн", 1},
		{"школа", "школьн", 2},
		{"сколь", "школьн", 2},
		{"школ", "школьн", 2},
		{"прикол", "школьн", 5},
		{"", "школьн", 6},
		{"школьн", "", 6},
		{"", "", 0},
	}

	for _, tt := range tests {
		t.Run(
			tt.a+"_"+tt.b, func(t *testing.T) {
				t.Parallel()

				if got := levenshtein(tt.a, tt.b); got != tt.want {
					t.Errorf("levenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
				}
			},
		)
	}
}
