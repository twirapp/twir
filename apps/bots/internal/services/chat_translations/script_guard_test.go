package chat_translations

import "testing"

func TestIsAmbiguousCyrillicDetection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		text     string
		detected string
		target   string
		want     bool
	}{
		{name: "russian misdetected as bulgarian", text: "чай попей", detected: "bg", target: "ru", want: true},
		{name: "russian with exclamation misdetected as bulgarian", text: "дайте токенов!", detected: "bg", target: "ru", want: true},
		{name: "short russian misdetected as bulgarian", text: "да не", detected: "bg", target: "ru", want: true},
		{name: "russian misdetected as serbian", text: "как дела братан", detected: "sr", target: "ru", want: true},
		{name: "russian misdetected as kazakh", text: "ёпта", detected: "kk", target: "ru", want: true},
		{name: "ukrainian with unique letters detected as ukrainian", text: "слава україні", detected: "uk", target: "ru", want: false},
		{name: "russian-only letter contradicts ukrainian detection", text: "привет ёж", detected: "uk", target: "ru", want: true},
		{name: "russian-only letter contradicts bulgarian detection", text: "эх ты", detected: "bg", target: "ru", want: true},
		{name: "serbian with unique letters detected as serbian", text: "ћао брате", detected: "sr", target: "ru", want: false},
		{name: "mixed cyrillic and latin with evidence for detected", text: "слава україні lol", detected: "uk", target: "ru", want: false},
		{name: "latin text is out of scope", text: "hello world", detected: "de", target: "ru", want: false},
		{name: "same language is handled earlier", text: "чай попей", detected: "ru", target: "ru", want: false},
		{name: "unknown detected language passes through", text: "чай попей", detected: "xx", target: "ru", want: false},
		{name: "non cyrillic target passes through", text: "чай попей", detected: "bg", target: "en", want: false},
		{name: "bulgarian without distinguishing letters is skipped by design", text: "добър ден", detected: "bg", target: "ru", want: true},
		{name: "kazakh with unique letters detected as kazakh", text: "сәлем қалайсың", detected: "kk", target: "ru", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isAmbiguousCyrillicDetection(tt.text, tt.detected, tt.target); got != tt.want {
				t.Fatalf(
					"isAmbiguousCyrillicDetection(%q, %q, %q) = %t, want %t",
					tt.text,
					tt.detected,
					tt.target,
					got,
					tt.want,
				)
			}
		})
	}
}
