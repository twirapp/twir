package chat_translations

import (
	"strings"
	"unicode"
)

// cyrillicAlphabets lists the lowercase letters of Cyrillic-script languages
// supported by the language detector. It is used to sanity-check detections:
// short Cyrillic messages are frequently misclassified between these
// languages (e.g. Russian detected as Bulgarian with confidence ~0.25).
var cyrillicAlphabets = map[string]string{
	"ru": "абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
	"uk": "абвгґдеєжзиіїйклмнопрстуфхцчшщьюя",
	"be": "абвгдеёжзійклмнопрстуўфхцчшыьэюя",
	"bg": "абвгдежзийклмнопрстуфхцчшщъьюя",
	"sr": "абвгдђежзијклљмнњопрстћуфхцчџш",
	"mk": "абвгдѓежзијклљмнњопрстќуфхцчџшѕ",
	"kk": "аәбвгғдеёжзийкқлмнңоөпрстуұүфхһцчшщъыіьэюя",
	"tg": "абвгғдеёжзиӣйкқлмнопрстуфхҳҷцчшъэюя",
	"mn": "абвгдеёжзийклмноөпрстуүфхцчшщъыьэюя",
	"ky": "абвгдеёжзийклмнңоөпрстуүфхцчшщъыьэюя",
}

// isAmbiguousCyrillicDetection reports whether a detected Cyrillic language
// cannot be distinguished from the target language by the letters actually
// present in the text (e.g. Russian "чай попей" misdetected as Bulgarian).
// It is true when the text has no letter unique to the detected alphabet, or
// has a letter impossible in it. Tradeoff: genuine messages in closely
// related languages without distinguishing letters are skipped too, but
// translating them produces near-identical text anyway.
func isAmbiguousCyrillicDetection(text, detectedLang, targetLang string) bool {
	if detectedLang == targetLang {
		return false
	}

	detectedAlphabet, ok := cyrillicAlphabets[detectedLang]
	if !ok {
		return false
	}
	targetAlphabet, ok := cyrillicAlphabets[targetLang]
	if !ok {
		return false
	}

	hasCyrillic := false
	evidenceForDetected := false

	for _, r := range strings.ToLower(text) {
		if !unicode.IsLetter(r) || !unicode.In(r, unicode.Cyrillic) {
			continue
		}

		hasCyrillic = true
		inDetected := strings.ContainsRune(detectedAlphabet, r)
		inTarget := strings.ContainsRune(targetAlphabet, r)

		if !inDetected && inTarget {
			// the letter cannot exist in the detected language,
			// but is valid in the target one
			return true
		}
		if inDetected && !inTarget {
			evidenceForDetected = true
		}
	}

	if !hasCyrillic {
		return false
	}

	return !evidenceForDetected
}
