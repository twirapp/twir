package moderation_helpers

import (
	"regexp"
)

func HasLink(regexp *regexp.Regexp, msg string) bool {
	var matches []string

	matches = regexp.FindAllString(msg, -1)

	return len(matches) > 0
}
