package moderation

import (
	"strings"

	"github.com/samber/lo"
)

func HasBlackListedWord(msg string, list []string) bool {
	msg = strings.ToLower(msg)

	return lo.SomeBy(list, func(i string) bool {
		return strings.Contains(msg, strings.ToLower(i))
	})
}
