package moderation

import "mvdan.cc/xurls/v2"

var linksMatcher = xurls.Relaxed()

func HasLink(msg string) bool {
	res := linksMatcher.FindAllString(msg, -1)
	return len(res) > 0
}
