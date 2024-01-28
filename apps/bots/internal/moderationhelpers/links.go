package moderationhelpers

func (c *ModerationHelpers) HasLink(msg string) bool {
	var matches []string

	matches = c.LinksWithSpaces.FindAllString(msg, -1)

	return len(matches) > 0
}
