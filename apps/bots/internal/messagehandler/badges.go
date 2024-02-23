package messagehandler

import (
	"strings"

	"github.com/satont/twir/libs/types/types/services/twitch"
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
