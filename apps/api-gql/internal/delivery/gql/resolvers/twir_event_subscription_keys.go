package resolvers

import (
	channelsservice "github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	twir_events "github.com/twirapp/twir/apps/api-gql/internal/services/twir-events"
)

func buildTwirEventSubscriptionKeys(identity channelsservice.ApiKeyChannelIdentity) []string {
	subscriptionKeys := []string{twir_events.CreateSubscribeKey(identity.InternalChannelID)}
	for _, target := range identity.ChatTargets {
		if target.Platform != "twitch" {
			continue
		}

		subscriptionKeys = append(subscriptionKeys, twir_events.CreateSubscribeKey(target.PlatformChannelID))
	}

	return subscriptionKeys
}
