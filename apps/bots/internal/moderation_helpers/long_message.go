package moderation_helpers

import (
	"unicode/utf8"
)

func IsTooLong(msg string, maxLength int) bool {
	return utf8.RuneCountInString(msg) > maxLength
}
