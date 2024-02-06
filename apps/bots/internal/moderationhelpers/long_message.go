package moderationhelpers

import (
	"unicode/utf8"
)

func (c *ModerationHelpers) IsTooLong(msg string, maxLength int) bool {
	return utf8.RuneCountInString(msg) > maxLength
}
