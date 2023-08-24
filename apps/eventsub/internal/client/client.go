package client

import (
	"context"
	"sync"
	"sync/atomic"

	eventsub_framework "github.com/dnsge/twitch-eventsub-framework"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/eventsub/internal/creds"
	"github.com/satont/twir/apps/eventsub/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"go.uber.org/zap"
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

	var channels []model.Channels
	err := services.Gorm.Where(
		`"isEnabled" = ? and "isBanned" = ?`,
		true,
		false,
	).Find(&channels).Error
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		err = subClient.SubscribeToNeededEvents(ctx, channel.ID)
		if err != nil {
			return nil, err
		}
	}

	return subClient, nil
}

type SubRequest struct {
	Version   string
	Condition map[string]string
}

func (c *SubClient) SubscribeToNeededEvents(ctx context.Context, userId string) error {
	channelCondition := map[string]string{
		"broadcaster_user_id": userId,
	}
	userCondition := map[string]string{
		"user_id": userId,
	}

	neededSubs := map[string]SubRequest{
		"channel.update": {
			Version:   "1",
			Condition: channelCondition,
		},
		"stream.online": {
			Version:   "1",
			Condition: channelCondition,
		},
		"stream.offline": {
			Version:   "1",
			Condition: channelCondition,
		},
		"user.update": {
			Condition: userCondition,
			Version:   "1",
		},
		"channel.follow": {
			Version: "2",
			Condition: map[string]string{
				"broadcaster_user_id": userId,
				"moderator_user_id":   userId,
			},
		},
		"channel.moderator.add": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.moderator.remove": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.channel_points_custom_reward_redemption.add": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.channel_points_custom_reward_redemption.update": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.poll.begin": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.poll.progress": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.poll.end": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.begin": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.lock": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.progress": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.end": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.ban": {
			Version:   "1",
			Condition: channelCondition,
		},
	}

	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		return err
	}

	subsReq, err := twitchClient.GetEventSubSubscriptions(
		&helix.EventSubSubscriptionsParams{
			UserID: userId,
		},
	)

	wg := &sync.WaitGroup{}
	wg.Add(len(neededSubs))

	var ops uint64

	for key, value := range neededSubs {
		go func(key string, value SubRequest) {
			defer wg.Done()

			existedSub, ok := lo.Find(
				subsReq.Data.EventSubSubscriptions, func(item helix.EventSubSubscription) bool {
					return item.Type == key &&
						(item.Condition.BroadcasterUserID == value.Condition["broadcaster_user_id"] ||
							item.Condition.UserID == value.Condition["user_id"])
				},
			)

			if ok && existedSub.Status == "enabled" && existedSub.Transport.Callback == c.callbackUrl {
				return
			}

			if ok {
				err = c.Unsubscribe(ctx, existedSub.ID)
				if err != nil {
					zap.S().Errorw("Failed to unsubscribe", "user_id", userId, "type", key, "error", err)
					return
				}
			}

			request := &eventsub_framework.SubRequest{
				Type:      key,
				Condition: value.Condition,
				Callback:  c.callbackUrl,
				Secret:    c.services.Config.TwitchClientSecret,
				Version:   value.Version,
			}
			if _, err := c.Subscribe(ctx, request); err != nil {
				zap.S().Error(err, key, userId)
				return
			}

			atomic.AddUint64(&ops, 1)
		}(key, value)
	}
	wg.Wait()

	zap.S().Infow("Subcribed to needed events", "user_id", userId, "madeRequests", ops)

	return nil
}
