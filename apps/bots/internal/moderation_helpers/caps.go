package moderation_helpers

import (
	"unicode"
	"unicode/utf8"
)

func IsTooMuchCaps(msg string, maxPercentage int) (bool, int) {
	capsCount := 0

	for _, s := range msg {
		if unicode.ToUpper(s) == s {
			capsCount++
		}
	}

	return capsCount*100 > maxPercentage*utf8.RuneCountInString(msg), capsCount
}
