package messagehandler

import (
	"strings"

	"github.com/twirapp/twir/libs/grpc/shared"
)

func createUserBadges(badges []*shared.ChatMessageBadge) []string {
	outer := make([]string, len(badges))

	for i, b := range badges {
		outer[i] = strings.ToUpper(b.GetSetId())

		if b.GetSetId() == "founder" {
			outer = append(outer, "SUBSCRIBER")
		}
	}

	return outer
}
