package moderationhelpers

import (
	"unicode"
	"unicode/utf8"
)

func (c *ModerationHelpers) IsTooMuchCaps(msg string, maxPercentage int) (bool, int) {
	capsCount := 0

	for _, s := range msg {
		if unicode.IsUpper(s) {
			capsCount++
		}
	}

	return capsCount*100 > maxPercentage*utf8.RuneCountInString(msg), capsCount
}
