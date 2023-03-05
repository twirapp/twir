package client

import (
	"context"
	"github.com/dnsge/twitch-eventsub-framework"
	"github.com/satont/tsuwari/apps/eventsub/internal/creds"
	"github.com/satont/tsuwari/apps/eventsub/internal/types"
)

type SubClient struct {
	*eventsub_framework.SubClient

	services    *types.Services
	callbackUrl string
}

func NewClient(
	ctx context.Context,
	services *types.Services,
	callBackUrl string,
) (*SubClient, error) {
	eventSubCreds := creds.NewCreds(ctx, services.Config, services.Grpc.Tokens)
	client := eventsub_framework.NewSubClient(eventSubCreds)

	subClient := &SubClient{
		SubClient:   client,
		services:    services,
		callbackUrl: callBackUrl,
	}

	if services.Config.AppEnv != "production" {
		subs, err := client.GetSubscriptions(ctx, eventsub_framework.StatusAny)
		if err != nil {
			return nil, err
		}

		for _, sub := range subs.Data {
			err = client.Unsubscribe(ctx, sub.ID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		subs, err := client.GetSubscriptions(ctx, eventsub_framework.StatusFailuresExceeded)
		if err != nil {
			return nil, err
		}

		for _, sub := range subs.Data {
			err = client.Unsubscribe(ctx, sub.ID)
			if err != nil {
				return nil, err
			}

			_, err = client.Subscribe(ctx, &eventsub_framework.SubRequest{
				Type:      sub.Type,
				Condition: sub.Condition,
				Callback:  callBackUrl,
				Secret:    services.Config.TwitchClientSecret,
			})
		}
	}

	return subClient, nil
}

func (c *SubClient) SubscribeToNeededEvents(ctx context.Context, userId string) error {
	channelCondition := map[string]string{
		"broadcaster_user_id": userId,
	}
	userCondition := map[string]string{
		"user_id": userId,
	}

	baseRequest := &eventsub_framework.SubRequest{
		Type: "",
		Condition: map[string]string{
			"broadcaster_user_id": userId,
		},
		Callback: c.callbackUrl,
		Secret:   c.services.Config.TwitchClientSecret,
	}

	baseRequest.Type = "channel.update"
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "stream.online"
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "stream.offline"
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "user.update"
	baseRequest.Condition = userCondition
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "channel.follow"
	baseRequest.Condition = map[string]string{
		"broadcaster_user_id": userId,
		"moderator_user_id":   userId,
	}
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "channel.moderator.add"
	baseRequest.Condition = channelCondition
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "channel.moderator.remove"
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "channel.channel_points_custom_reward_redemption.add"
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	baseRequest.Type = "channel.channel_points_custom_reward_redemption.update"
	if _, err := c.Subscribe(ctx, baseRequest); err != nil {
		return err
	}

	return nil
}
