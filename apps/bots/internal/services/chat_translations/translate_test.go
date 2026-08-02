package chat_translations

import "testing"

func TestPrepareTextForTranslation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "url only",
			text: "https://www.youtube.com/watch?v=pM9rYSN7WWY&list=PLrm7ivpQ4YpWl-OO_JIjU1zQjeWny7wqP&index=19",
			want: "",
		},
		{
			name: "text followed by url",
			text: "check this https://example.com/watch?v=123&list=456",
			want: "check this",
		},
		{
			name: "www url between text",
			text: "before www.example.com/path after",
			want: "before after",
		},
		{
			name: "normalizes whitespace",
			text: "  hello\n\tworld  ",
			want: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := prepareTextForTranslation(tt.text); got != tt.want {
				t.Fatalf("prepareTextForTranslation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHasTranslatableText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "natural text", text: "тетс триггера перевода", want: true},
		{name: "numbers and punctuation", text: "12345 !!!", want: false},
		{name: "too few letters", text: "ok 12345", want: false},
		{name: "non latin letters", text: "你好世界", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := hasTranslatableText(tt.text); got != tt.want {
				t.Fatalf("hasTranslatableText() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestTranslationsEquivalent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		source     string
		translated string
		want       bool
	}{
		{name: "identical text", source: "hello world", translated: "hello world", want: true},
		{name: "case and punctuation shuffle", source: "чай попей", translated: "чай Попай", want: true},
		{name: "single letter suffix", source: "да не", translated: "да нет", want: true},
		{name: "different punctuation only", source: "привет, как дела?", translated: "Привет как дела", want: true},
		{name: "synonym swap is not equivalent", source: "дайте токенов!", translated: "Дайте мне жетон!", want: false},
		{name: "real translation", source: "hello my friend", translated: "привет мой друг", want: false},
		{name: "empty source", source: "!!!", translated: "???", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := translationsEquivalent(tt.source, tt.translated); got != tt.want {
				t.Fatalf(
					"translationsEquivalent(%q, %q) = %t, want %t",
					tt.source,
					tt.translated,
					got,
					tt.want,
				)
			}
		})
	}
}
