package eventsub_bindings

type ConditionChannelBan struct {
	// The broadcaster user ID for the channel you want to get ban notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelSubscribe struct {
	// The broadcaster user ID for the channel you want to get subscribe notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelSubscriptionEnd struct {
	// The broadcaster user ID for the channel you want to get subscription end notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelSubscriptionGift struct {
	// The broadcaster user ID for the channel you want to get subscription gift notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelSubscriptionMessage struct {
	// The broadcaster user ID for the channel you want to get resubscription chat message notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelCheer struct {
	// The broadcaster user ID for the channel you want to get cheer notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelUpdate struct {
	// The broadcaster user ID for the channel you want to get updates for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelFollow struct {
	// The broadcaster user ID for the channel you want to get follow notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelUnban struct {
	// The broadcaster user ID for the channel you want to get unban notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelRaid struct {
	// The broadcaster user ID that created the channel raid you want to get notifications for.
	// Use this parameter if you want to know when a specific broadcaster raids another broadcaster.
	//
	// This field is optional if ToBroadcasterUserID is set.
	FromBroadcasterUserID string `json:"from_broadcaster_user_id,omitempty"`
	// The broadcaster user ID that received the channel raid you want to get notifications for.
	// Use this parameter if you want to know when a specific broadcaster is raided by another broadcaster.
	//
	// This field is optional if FromBroadcasterUserID is set.
	ToBroadcasterUserID string `json:"to_broadcaster_user_id,omitempty"`
}

type ConditionChannelModeratorAdd struct {
	// The broadcaster user ID for the channel you want to get moderator addition notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelModeratorRemove struct {
	// The broadcaster user ID for the channel you want to get moderator removal notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPointsRewardAdd struct {
	// The broadcaster user ID for the channel you want to receive channel points custom reward add notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPointsRewardUpdate struct {
	// The broadcaster user ID for the channel you want to receive channel points custom reward update notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// Optional. Specify a reward id to only receive notifications for a specific reward.
	RewardID string `json:"reward_id,omitempty"`
}

type ConditionChannelPointsRewardRemove struct {
	// The broadcaster user ID for the channel you want to receive channel points custom reward remove notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// Optional. Specify a reward id to only receive notifications for a specific reward.
	RewardID string `json:"reward_id,omitempty"`
}

type ConditionChannelPointsRewardRedemptionAdd struct {
	// The broadcaster user ID for the channel you want to receive channel points custom reward redemption add notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// Optional. Specify a reward id to only receive notifications for a specific reward.
	RewardID string `json:"reward_id,omitempty"`
}

type ConditionChannelPointsRewardRedemptionUpdate struct {
	// The broadcaster user ID for the channel you want to receive channel points custom reward redemption update notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// Optional. Specify a reward id to only receive notifications for a specific reward.
	RewardID string `json:"reward_id,omitempty"`
}

type ConditionChannelPollBegin struct {
	// The broadcaster user ID of the channel for which "poll begin" notifications will be received.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPollProgress struct {
	// The broadcaster user ID of the channel for which "poll progress" notifications will be received.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPollEnd struct {
	// The broadcaster user ID of the channel for which "poll end" notifications will be received.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPredictionBegin struct {
	// The broadcaster user ID of the channel for which "prediction begin" notifications will be received.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPredictionLock struct {
	// The broadcaster user ID of the channel for which "prediction progress" notifications will be received.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionChannelPredictionEnd struct {
	// The broadcaster user ID of the channel for which "prediction end" notifications will be received.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionDropEntitlementGrant struct {
	// The organization ID of the organization that owns the game on the developer portal.
	OrganizationID string `json:"organization_id"`
	// Optional. The category (or game) ID of the game for which entitlement notifications will be received.
	CategoryID string `json:"category_id,omitempty"`
	// Optional. The campaign ID for a specific campaign for which entitlement notifications will be received.
	CampaignID string `json:"campaign_id,omitempty"`
}

type ConditionExtensionBitsTransactionCreate struct {
	// The client ID of the extension.
	ExtensionClientID string `json:"extension_client_id"`
}

type ConditionGoals struct {
	// The ID of the broadcaster to get notified about. The ID must match the user_id in the OAuth access token.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionHypeTrainBegin struct {
	// The broadcaster user ID for the channel you want to Hype Train begin notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionHypeTrainProgress struct {
	// The broadcaster user ID for the channel you want to Hype Train progress notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionHypeTrainEnd struct {
	// The broadcaster user ID for the channel you want to Hype Train end notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionStreamOnline struct {
	// The broadcaster user ID you want to get stream online notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionStreamOffline struct {
	// The broadcaster user ID you want to get stream offline notifications for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type ConditionUserAuthorizationGrant struct {
	// Your application’s client id. The provided client_id must match the client id in the application access token.
	ClientID string `json:"client_id"`
}

type ConditionUserAuthorizationRevoke struct {
	// Your application’s client id. The provided client_id must match the client id in the application access token.
	ClientID string `json:"client_id"`
}

type ConditionUserUpdate struct {
	// The user ID for the user you want update notifications for.
	UserID string `json:"user_id"`
}
