package moderationhelpers

import "regexp"

var (
	clipsRegex1 = regexp.MustCompile(`(?i).*(clips\.twitch\.tv/)(\w+)`)
	clipsRegex2 = regexp.MustCompile(`(?i).*(www\.twitch\.tv/\w+/clip/)(\w+)`)
)

func (c *ModerationHelpers) HasLink(msg string, checkClips bool) bool {
	matches := c.LinksWithSpaces.FindAllString(msg, -1)

	if len(matches) == 0 {
		return false
	}

	if !checkClips {
		var clipsLength int
		clipsLength += len(clipsRegex1.FindAllString(msg, -1))
		clipsLength += len(clipsRegex2.FindAllString(msg, -1))
		if clipsLength == len(matches) {
			return false
		}
	}

	return true
}
