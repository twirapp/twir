package moderationhelpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/twirapp/twir/apps/bots/pkg/tlds"
)

type ModerationHelpers struct {
	LinksWithSpaces *regexp.Regexp
}

func New(tlds *tlds.TLDS) *ModerationHelpers {
	return &ModerationHelpers{
		LinksWithSpaces: regexp.MustCompile(
			fmt.Sprintf(
				`(www)? ??\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\. ?(%s)\b`,
				strings.Join(tlds.List, "|"),
			),
		),
	}
}
