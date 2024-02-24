package messagehandler

import (
	"strings"

	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func createUserBadges(badges []twitch.ChatMessageBadge) []string {
	outer := make([]string, len(badges))

	for i, b := range badges {
		outer[i] = strings.ToUpper(b.SetId)

		if b.SetId == "founder" {
			outer = append(outer, "SUBSCRIBER")
		}
	}

	return outer
}
