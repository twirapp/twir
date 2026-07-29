package manager

import (
	"errors"

	"github.com/kvizyx/twitchy/eventsub"
)

func (c *Manager) getConditionForTopic(eventType eventsub.EventType, broadcasterID string, botID string) (eventsub.Condition, error) {
	switch eventType {
	case eventsub.EventTypeAutomodMessageHold:
		return eventsub.AutomodMessageHoldCondition{BroadcasterUserId: broadcasterID, ModeratorUserId: botID}, nil
	case eventsub.EventTypeUserAuthorizationRevoke:
		return eventsub.UserAuthorizationRevokeCondition{ClientId: c.config.TwitchClientId}, nil
	case eventsub.EventTypeChannelFollow:
		return eventsub.ChannelFollowCondition{BroadcasterUserId: broadcasterID, ModeratorUserId: botID}, nil
	case eventsub.EventTypeChannelBan:
		return eventsub.ChannelBanCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelUnban:
		return eventsub.ChannelUnbanCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelChatClear:
		return eventsub.ChannelChatClearCondition{BroadcasterUserId: broadcasterID, UserId: botID}, nil
	case eventsub.EventTypeChannelChatClearUserMessages:
		return eventsub.ChannelChatClearUserMessagesCondition{BroadcasterUserId: broadcasterID, UserId: botID}, nil
	case eventsub.EventTypeChannelChatMessage:
		return eventsub.ChannelChatMessageCondition{BroadcasterUserId: broadcasterID, UserId: botID}, nil
	case eventsub.EventTypeChannelChatNotification:
		return eventsub.ChannelChatNotificationCondition{BroadcasterUserId: broadcasterID, UserId: botID}, nil
	case eventsub.EventTypeChannelModeratorAdd:
		return eventsub.ChannelModeratorAddCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelModeratorRemove:
		return eventsub.ChannelModeratorRemoveCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPollBegin:
		return eventsub.ChannelPollBeginCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPollProgress:
		return eventsub.ChannelPollProgressCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPollEnd:
		return eventsub.ChannelPollEndCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPredictionBegin:
		return eventsub.ChannelPredictionBeginCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPredictionProgress:
		return eventsub.ChannelPredictionProgressCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPredictionLock:
		return eventsub.ChannelPredictionLockCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPredictionEnd:
		return eventsub.ChannelPredictionEndCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelRaid:
		return eventsub.ChannelRaidCondition{ToBroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPointsCustomRewardRedemptionAdd:
		return eventsub.ChannelPointsCustomRewardRedemptionAddCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPointsCustomRewardRedemptionUpdate:
		return eventsub.ChannelPointsCustomRewardRedemptionUpdateCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPointsAutomaticRewardRedemptionAdd:
		return eventsub.ChannelPointsAutomaticRewardRedemptionAddCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPointsRewardAdd:
		return eventsub.ChannelPointsCustomRewardAddCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPointsRewardUpdate:
		return eventsub.ChannelPointsCustomRewardUpdateCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelPointsRewardRemove:
		return eventsub.ChannelPointsCustomRewardRemoveCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeStreamOffline:
		return eventsub.StreamOfflineCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeStreamOnline:
		return eventsub.StreamOnlineCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelSubscribe:
		return eventsub.ChannelSubscribeCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelSubscriptionEnd:
		return eventsub.ChannelSubscriptionEndCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelSubscriptionMessage:
		return eventsub.ChannelSubscriptionMessageCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelSubscriptionGift:
		return eventsub.ChannelSubscriptionGiftCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelUnbanRequestCreate:
		return eventsub.ChannelUnbanRequestCreateCondition{BroadcasterUserId: broadcasterID, ModeratorUserId: botID}, nil
	case eventsub.EventTypeChannelUnbanRequestResolve:
		return eventsub.ChannelUnbanRequestResolveCondition{BroadcasterUserId: broadcasterID, ModeratorUserId: botID}, nil
	case eventsub.EventTypeUserUpdate:
		return eventsub.UserUpdateCondition{UserId: broadcasterID}, nil
	case eventsub.EventTypeChannelVipAdd:
		return eventsub.ChannelVIPAddCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelVipRemove:
		return eventsub.ChannelVIPRemoveCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelMessageDelete:
		return eventsub.ChannelChatMessageDeleteCondition{BroadcasterUserId: broadcasterID, UserId: botID}, nil
	case eventsub.EventTypeChannelUpdate:
		return eventsub.ChannelUpdateCondition{BroadcasterUserId: broadcasterID}, nil
	case eventsub.EventTypeChannelModerate:
		return eventsub.ChannelModerateV2Condition{BroadcasterUserId: broadcasterID, ModeratorUserId: botID}, nil
	default:
		return nil, errors.New("unsupported event type for topic")
	}
}
