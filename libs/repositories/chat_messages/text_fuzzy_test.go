package chat_messages

import "testing"

func TestNewTextFuzzyFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		phrase  string
		want    *TextFuzzyFilter
		wantNil bool
	}{
		{name: "single word", phrase: "школьн", want: &TextFuzzyFilter{Phrase: "школьн", Length: 6, MaxDistance: 1}},
		{name: "normalized case and spaces", phrase: "  Школьн ", want: &TextFuzzyFilter{Phrase: "школьн", Length: 6, MaxDistance: 1}},
		{name: "long phrase two edits", phrase: "школьницы", want: &TextFuzzyFilter{Phrase: "школьницы", Length: 9, MaxDistance: 2}},
		{name: "punctuation stripped", phrase: "школьн!", want: &TextFuzzyFilter{Phrase: "школьн", Length: 6, MaxDistance: 1}},
		{name: "short phrase", phrase: "школ", wantNil: true},
		{name: "multi word", phrase: "голые школьницы", wantNil: true},
		{name: "empty", phrase: "", wantNil: true},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				got := NewTextFuzzyFilter(tt.phrase)
				if tt.wantNil {
					if got != nil {
						t.Errorf("NewTextFuzzyFilter(%q) = %+v, want nil", tt.phrase, got)
					}
					return
				}

				if got == nil || *got != *tt.want {
					t.Errorf("NewTextFuzzyFilter(%q) = %+v, want %+v", tt.phrase, got, tt.want)
				}
			},
		)
	}
}
