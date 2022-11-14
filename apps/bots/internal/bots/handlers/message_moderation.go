package handlers

import (
	"regexp"

	"mvdan.cc/xurls/v2"
)

var clipsRegexps = [2]*regexp.Regexp{
	regexp.MustCompile(`.*(clips.twitch.tv\/)(\w+)`),
	regexp.MustCompile(`.*(www.twitch.tv\/\w+\/clip\/)(\w+)`),
}

var linksMatcher = xurls.Relaxed()

func (c *Handlers) handleModerate() {
}
