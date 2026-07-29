package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/avast/retry-go/v4"
	"github.com/kvizyx/twitchy/eventsub"
	"github.com/kvizyx/twitchy/helix"
)

func (c *Manager) SubscribeWithLimits(
	ctx context.Context,
	eventType eventsub.EventType,
	eventTransport eventsub.Transport,
	eventVersion string,
	broadcasterId string,
	botId string,
) error {
	condition, err := c.getConditionForTopic(eventType, broadcasterId, botId)
	if err != nil {
		return err
	}

	conditionBytes, err := json.Marshal(&condition)
	if err != nil {
		return err
	}

	transportBytes, err := json.Marshal(&eventTransport)
	if err != nil {
		return err
	}

	conditionMap := helix.EventSubCondition{}
	if err := json.Unmarshal(conditionBytes, &conditionMap); err != nil {
		return err
	}

	transport := helix.EventSubTransport{}
	if err := json.Unmarshal(transportBytes, &transport); err != nil {
		return err
	}

	requestData := helix.CreateEventSubSubscriptionRequest{
		Type:      eventType.String(),
		Version:   eventVersion,
		Condition: conditionMap,
		Transport: transport,
	}

	var twitchClient *helix.Client
	switch eventTransport.(type) {
	case eventsub.ConduitTransport:
		client, err := c.newAppTwitchClient(ctx)
		if err != nil {
			return err
		}
		twitchClient = client
	case eventsub.WebsocketTransport:
		client, err := c.newBotTwitchClient(ctx, botId)
		if err != nil {
			return err
		}
		twitchClient = client
	case eventsub.WebhookTransport:
		client, err := c.newAppTwitchClient(ctx)
		if err != nil {
			return err
		}
		twitchClient = client
	default:
		return fmt.Errorf("unsupported EventSub transport %T", eventTransport)
	}

	return retry.Do(
		func() error {
			_, err := twitchClient.EventSub.CreateEventSubSubscription(ctx, requestData)
			return err
		},
		retry.Attempts(0),
		retry.RetryIf(
			func(err error) bool {
				if strings.Contains(err.Error(), "maximum subscriptions with type and condition exceeded") {
					return false
				}
				var rateLimitErr *helix.RateLimitError
				return errors.As(err, &rateLimitErr)
			},
		),
	)
}
