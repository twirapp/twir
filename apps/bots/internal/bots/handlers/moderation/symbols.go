package moderation

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var symbolsRegexp = regexp.MustCompile(`[^\s\\u0500-\\u052F\\u0400-\\u04FF\w]+`)

func IsToMuchSymbols(msg string, maxPercentage int) bool {
	msg = strings.ReplaceAll(msg, " ", "")
	matches := symbolsRegexp.FindAllString(msg, -1)
	matchesCount := 0

	for _, v := range matches {
		matchesCount += utf8.RuneCountInString(v)
	}

	return matchesCount*100 >= maxPercentage*utf8.RuneCountInString(msg)
}
