package seventv

import (
	"regexp"
)

var emoteRegex = regexp.MustCompile(`((cdn.)?7tv.app/emotes/)(?P<id>.{26})`)

func FindEmoteIdInInput(input string) string {
	var result string
	groupNames := emoteRegex.SubexpNames()
	for _, match := range emoteRegex.FindAllStringSubmatch(input, -1) {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]
			if name == "id" {
				result = group
				break
			}
		}
	}

	return result
}
