package moderation_helpers

import (
	"regexp"
	"strings"

	"github.com/samber/lo"
)

func HasDeniedWord(msg string, list []string) bool {
	msg = strings.ToLower(msg)

	return lo.SomeBy(
		list,
		func(item string) bool {
			r, err := regexp.Compile(item)
			if err == nil {
				return r.MatchString(msg)
			}

			return strings.Contains(msg, strings.ToLower(item))
		},
	)
}
