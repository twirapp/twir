package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/imroc/req/v3"
	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/tokens"
)

type ErrRateLimit struct {
	Err        error
	RetryAfter time.Duration
}

func (e *ErrRateLimit) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return "rate limit exceeded"
}

func (c *Manager) SubscribeWithLimits(
	ctx context.Context,
	eventType eventsub.EventType,
	eventTransport eventsub.Transport,
	eventVersion string,
	channelId,
	botId string,
) error {
	condition, err := c.getConditionForTopic(eventType, channelId, botId)
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

	conditionMap := map[string]any{}
	if err := json.Unmarshal(conditionBytes, &conditionMap); err != nil {
		return err
	}

	transportMap := map[string]any{}
	if err := json.Unmarshal(transportBytes, &transportMap); err != nil {
		return err
	}

	requestData := map[string]any{
		"type":      eventType.String(),
		"version":   eventVersion,
		"condition": conditionMap,
		"transport": transportMap,
	}

	requestBytes, err := json.Marshal(&requestData)
	if err != nil {
		return err
	}

	reqBuilder := req.R().
		SetContext(ctx).
		SetHeader(
			"Client-id",
			c.config.TwitchClientId,
		)

	switch eventTransport.(type) {
	case eventsub.ConduitTransport:
		appToken, err := c.twirBus.Tokens.RequestAppToken.Request(
			ctx,
			struct{}{},
		)
		if err != nil {
			return err
		}

		reqBuilder = reqBuilder.
			SetBearerAuthToken(appToken.Data.AccessToken)
	case eventsub.WebsocketTransport:
		botToken, err := c.twirBus.Tokens.RequestBotToken.Request(
			ctx,
			tokens.GetBotTokenRequest{
				BotId: botId,
			},
		)
		if err != nil {
			return err
		}

		reqBuilder = reqBuilder.
			SetBearerAuthToken(botToken.Data.AccessToken)
	case eventsub.WebhookTransport:
		appToken, err := c.twirBus.Tokens.RequestAppToken.Request(
			ctx,
			struct{}{},
		)
		if err != nil {
			return err
		}

		reqBuilder = reqBuilder.
			SetBearerAuthToken(appToken.Data.AccessToken)
	}

	err = retry.Do(
		func() error {
			resp, err := reqBuilder.SetBodyJsonBytes(requestBytes).Post("https://api.twitch.tv/helix/eventsub/subscriptions")
			if err != nil {
				return err
			}

			if resp.IsErrorState() {
				if resp.StatusCode == 429 && !strings.Contains(
					resp.String(),
					"maximum subscriptions with type and condition exceeded",
				) {
					resetTimeStr := resp.Header.Get("ratelimit-reset")
					if resetTimeStr == "" {
						return fmt.Errorf("ratelimit-reset header not found")
					}
					resetTimeUnix, err := strconv.ParseInt(resetTimeStr, 10, 64)
					if err != nil {
						return fmt.Errorf("failed to parse ratelimit-reset: %v", err)
					}
					resetTime := time.Unix(resetTimeUnix, 0)
					retryAfter := time.Until(resetTime)
					if retryAfter <= 0 {
						retryAfter = time.Second
					}

					return &ErrRateLimit{
						Err:        fmt.Errorf("rate limit exceeded"),
						RetryAfter: retryAfter,
					}
				}

				return errors.New(resp.String())
			}

			return nil
		},
		retry.Attempts(0),
		retry.DelayType(
			func(n uint, err error, config *retry.Config) time.Duration {
				var rateLimitErr *ErrRateLimit
				if errors.As(err, &rateLimitErr) {
					return rateLimitErr.RetryAfter
				}

				return retry.BackOffDelay(n, err, config)
			},
		),
		retry.RetryIf(
			func(err error) bool {
				var rateLimitErr *ErrRateLimit
				if errors.As(err, &rateLimitErr) {
					// If the error is a rate limit error, we should retry
					return true
				}

				if errors.Is(err, context.DeadlineExceeded) {
					return true
				}

				if strings.Contains(err.Error(), "context deadline exceeded") {
					return true
				}

				return false
			},
		),
	)

	return err
}

func (c *Manager) getConditionForTopic(
	eventType eventsub.EventType,
	channelId, botId string,
) (eventsub.Condition, error) {
	switch eventType {
	case eventsub.EventTypeAutomodMessageHold:
		return eventsub.AutomodMessageHoldCondition{
			BroadcasterUserId: channelId,
			ModeratorUserId:   botId,
		}, nil
	case eventsub.EventTypeUserAuthorizationRevoke:
		return eventsub.UserAuthorizationRevokeCondition{
			ClientId: c.config.TwitchClientId,
		}, nil
	case eventsub.EventTypeChannelFollow:
		return eventsub.ChannelFollowCondition{
			BroadcasterUserId: channelId,
			ModeratorUserId:   botId,
		}, nil
	case eventsub.EventTypeChannelBan:
		return eventsub.ChannelBanCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelUnban:
		return eventsub.ChannelUnbanCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelChatClear:
		return eventsub.ChannelChatClearCondition{
			BroadcasterUserId: channelId,
			UserId:            botId,
		}, nil
	case eventsub.EventTypeChannelChatClearUserMessages:
		return eventsub.ChannelChatClearUserMessagesCondition{
			BroadcasterUserId: channelId,
			UserId:            botId,
		}, nil
	case eventsub.EventTypeChannelChatMessage:
		return eventsub.ChannelChatMessageCondition{
			BroadcasterUserId: channelId,
			UserId:            botId,
		}, nil
	case eventsub.EventTypeChannelChatNotification:
		return eventsub.ChannelChatNotificationCondition{
			BroadcasterUserId: channelId,
			UserId:            botId,
		}, nil
	case eventsub.EventTypeChannelModeratorAdd:
		return eventsub.ChannelModeratorAddCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelModeratorRemove:
		return eventsub.ChannelModeratorRemoveCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPollBegin:
		return eventsub.ChannelPollBeginCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPollProgress:
		return eventsub.ChannelPollProgressCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPollEnd:
		return eventsub.ChannelPollEndCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPredictionBegin:
		return eventsub.ChannelPredictionBeginCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPredictionProgress:
		return eventsub.ChannelPredictionProgressCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPredictionLock:
		return eventsub.ChannelPredictionLockCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPredictionEnd:
		return eventsub.ChannelPredictionEndCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelRaid:
		return eventsub.ChannelRaidCondition{
			ToBroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPointsCustomRewardRedemptionAdd:
		return eventsub.ChannelPointsCustomRewardRedemptionAddCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPointsCustomRewardRedemptionUpdate:
		return eventsub.ChannelPointsCustomRewardRedemptionUpdateCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPointsAutomaticRewardRedemptionAdd:
		return eventsub.ChannelPointsAutomaticRewardRedemptionAddCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPointsRewardAdd:
		return eventsub.ChannelPointsCustomRewardAddCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPointsRewardUpdate:
		return eventsub.ChannelPointsCustomRewardUpdateCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelPointsRewardRemove:
		return eventsub.ChannelPointsCustomRewardRemoveCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeStreamOffline:
		return eventsub.StreamOfflineCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeStreamOnline:
		return eventsub.StreamOnlineCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelSubscribe:
		return eventsub.ChannelSubscribeCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelSubscriptionEnd:
		return eventsub.ChannelSubscriptionEndCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelSubscriptionMessage:
		return eventsub.ChannelSubscriptionMessageCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelSubscriptionGift:
		return eventsub.ChannelSubscriptionGiftCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelUnbanRequestCreate:
		return eventsub.ChannelUnbanRequestCreateCondition{
			BroadcasterUserId: channelId,
			ModeratorUserId:   botId,
		}, nil
	case eventsub.EventTypeChannelUnbanRequestResolve:
		return eventsub.ChannelUnbanRequestResolveCondition{
			BroadcasterUserId: channelId,
			ModeratorUserId:   botId,
		}, nil
	case eventsub.EventTypeUserUpdate:
		return eventsub.UserUpdateCondition{
			UserId: channelId,
		}, nil
	case eventsub.EventTypeChannelVipAdd:
		return eventsub.ChannelVIPAddCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelVipRemove:
		return eventsub.ChannelVIPRemoveCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelMessageDelete:
		return eventsub.ChannelChatMessageDeleteCondition{
			BroadcasterUserId: channelId,
			UserId:            botId,
		}, nil
	case eventsub.EventTypeChannelUpdate:
		return eventsub.ChannelUpdateCondition{
			BroadcasterUserId: channelId,
		}, nil
	case eventsub.EventTypeChannelModerate:
		return eventsub.ChannelModerateV2Condition{
			BroadcasterUserId: channelId,
			ModeratorUserId:   botId,
		}, nil
	default:
		return nil, errors.New("unsupported event type for topic")
	}
}
