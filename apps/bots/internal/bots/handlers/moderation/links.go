package moderation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/satont/tsuwari/apps/bots/pkg/tlds"
)

// [a-zA-Z0-9]+([a-zA-Z0-9-]+)?\\.(${tlds.join('|')})
var linksWithSpaces = regexp.MustCompile(
	fmt.Sprintf(
		`(www)? ??\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\. ?(%s)\b`,
		strings.Join(tlds.TLDS, "|"),
	),
)

var linksWithoutSpaces = regexp.MustCompile(
	fmt.Sprintf(
		`[a-zA-Z0-9]+([a-zA-Z0-9-]+)?\.(%s)\b`,
		strings.Join(tlds.TLDS, "|"),
	),
)

func HasLink(msg string, withSpaces bool) bool {
	var matches []string

	if withSpaces {
		matches = linksWithSpaces.FindAllString(msg, -1)
	} else {
		matches = linksWithoutSpaces.FindAllString(msg, -1)
	}

	return len(matches) > 0
}
