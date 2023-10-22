package moderation_helpers

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var symbolsRegexp = regexp.MustCompile(`[^\p{L}0-9\s]+`)

// IsToMuchSymbols checks if message contains more than maxPercentage of symbols
// Take care for one emoji it counts as 2 symbols
func IsToMuchSymbols(msg string, maxPercentage int) (bool, int) {
	msg = strings.ReplaceAll(msg, " ", "")
	matches := symbolsRegexp.FindAllString(msg, -1)
	matchesCount := 0

	for _, v := range matches {
		matchesCount += utf8.RuneCountInString(v)
	}

	return matchesCount*100 >= maxPercentage*utf8.RuneCountInString(msg), matchesCount
}
