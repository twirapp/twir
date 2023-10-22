package moderation_helpers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasLink(t *testing.T) {
	linksR, linksWithSpaces := BuildLinksModerationRegexps([]string{"ru", "com", "tv"})

	table := []struct {
		msg           string
		want          bool
		wantOnlyClips bool
		regular       *regexp.Regexp
	}{
		{"https://www.google.com", true, false, linksR},
		{"https://www.google.by", false, false, linksR},
		{"https://www.google . com", false, false, linksR},

		{"https://www.google . com", true, false, linksWithSpaces},
		{"https://www.google.com", true, false, linksWithSpaces},

		// {"https://www.twitch.tv/satont/clip/qwe", true, true, linksWithSpaces},
		// {
		// 	"https://google.com and https://www.twitch.tv/satont/clip/qwe",
		// 	true,
		// 	false,
		// 	linksWithSpaces,
		// },
		//
		// {"https://clips.twitch.tv/qwe", true, true, linksR},
		// {"https://google.com and https://clips.twitch.tv/asd", true, false, linksWithSpaces},
	}

	for _, tt := range table {
		t.Run(
			tt.msg,
			func(t *testing.T) {
				got := HasLink(tt.regular, tt.msg)

				assert.Equal(t, tt.want, got, tt.msg)
			},
		)
	}
}
