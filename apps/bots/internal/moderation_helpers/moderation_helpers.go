package moderation_helpers

import (
	"fmt"
	"regexp"
	"strings"
)

func BuildLinksModerationRegexps(tldsList []string) (*regexp.Regexp, *regexp.Regexp) {
	linksWithoutSpaces := regexp.MustCompile(
		fmt.Sprintf(
			`[a-zA-Z0-9]+([a-zA-Z0-9-]+)?\.(%s)\b`,
			strings.Join(tldsList, "|"),
		),
	)

	linksWithSpaces := regexp.MustCompile(
		fmt.Sprintf(
			`(www)? ??\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\. ?(%s)\b`,
			strings.Join(tldsList, "|"),
		),
	)

	return linksWithoutSpaces, linksWithSpaces
}
