package chat_alerts

import (
	"cmp"
	"slices"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *ChatAlerts) takeCountedSample(
	target int,
	messages []model.ChatAlertsCountedMessage,
) string {
	if len(messages) == 0 {
		return ""
	}

	slices.SortFunc(
		messages, func(a, b model.ChatAlertsCountedMessage) int {
			return cmp.Compare(a.Count, b.Count)
		},
	)

	var lastMatch model.ChatAlertsCountedMessage
	for _, m := range messages {
		if m.Count <= target {
			lastMatch = m
		}
	}

	groupedMatched := lo.GroupBy(
		messages,
		func(m model.ChatAlertsCountedMessage) int {
			return m.Count
		},
	)
	sample := lo.Sample(groupedMatched[lastMatch.Count]).Text

	return sample
}
