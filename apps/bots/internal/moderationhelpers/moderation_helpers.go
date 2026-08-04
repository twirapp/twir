package moderationhelpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/twirapp/twir/apps/bots/pkg/tlds"
)

type Opts struct {
	Tlds *tlds.TLDS
}

type ModerationHelpers struct {
	LinksWithSpaces *regexp.Regexp
}

func New(opts Opts) *ModerationHelpers {
	return &ModerationHelpers{
		LinksWithSpaces: regexp.MustCompile(
			fmt.Sprintf(
				`(www)? ??\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\. ?(%s)\b`,
				strings.Join(opts.Tlds.List, "|"),
			),
		),
	}
}
