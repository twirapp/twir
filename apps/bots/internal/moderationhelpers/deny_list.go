package moderationhelpers

import (
	"regexp"
	"strings"
)

const denyUnicodeAwareBoundaryPrefix = `(?i)(?:^|\s|[^\p{L}\p{N}])`
const denyUnicodeAwareBoundarySuffix = `(?:$|\s|[^\p{L}\p{N}])`

func (c *ModerationHelpers) HasDeniedWord(msg string, list []string) bool {
	msg = strings.ToLower(msg)

	for _, item := range list {
		if item == "" {
			continue
		}

		r, err := regexp.Compile(item)
		if err == nil {
			matched := r.MatchString(msg)
			if matched {
				return true
			}
		}

		wordRg := regexp.MustCompile(denyUnicodeAwareBoundaryPrefix + regexp.QuoteMeta(strings.ToLower(item)) + denyUnicodeAwareBoundarySuffix)
		matched := wordRg.MatchString(msg)

		if matched {
			return true
		}
	}

	return false
}
