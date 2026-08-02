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
