package eventsub_bindings

type EventChannelBan struct {
	// The user ID for the user who was banned on the specified channel.
	UserID string `json:"user_id"`
	// The user login for the user who was banned on the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user who was banned on the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the issuer of the ban.
	ModeratorUserID string `json:"moderator_user_id"`
	// The user login of the issuer of the ban.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user name of the issuer of the ban.
	ModeratorUserName string `json:"moderator_user_name"`
	// The reason behind the ban.
	Reason string `json:"reason"`
	// Will be null if permanent ban. If it is a timeout, this field shows when the timeout will end.
	EndsAt string `json:"ends_at"`
	// Indicates whether the ban is permanent (true) or a timeout (false). If true, ends_at will be null.
	IsPermanent bool `json:"is_permanent"`
}

type EventChannelSubscribe struct {
	// The user ID for the user who subscribed to the specified channel.
	UserID string `json:"user_id"`
	// The user login for the user who subscribed to the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user who subscribed to the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The tier of the subscription. Valid values are 1000, 2000, and 3000.
	Tier string `json:"tier"`
	// Whether the subscription is a gift.
	IsGift bool `json:"is_gift"`
}

type EventChannelCheer struct {
	// Whether the user cheered anonymously or not.
	IsAnonymous bool `json:"is_anonymous"`
	// The user ID for the user who cheered on the specified channel. This is null if is_anonymous is true.
	UserID string `json:"user_id"`
	// The user login for the user who cheered on the specified channel. This is null if is_anonymous is true.
	UserLogin string `json:"user_login"`
	// The user display name for the user who cheered on the specified channel. This is null if is_anonymous is true.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The message sent with the cheer.
	Message string `json:"message"`
	// The number of bits cheered.
	Bits int `json:"bits"`
}

type EventChannelUpdate struct {
	// The broadcaster’s user ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The channel’s stream title.
	Title string `json:"title"`
	// The channel’s broadcast language.
	Language string `json:"language"`
	// The channel’s category ID.
	CategoryID string `json:"category_id"`
	// The category name.
	CategoryName string `json:"category_name"`
	// A boolean identifying whether the channel is flagged as mature. Valid values are true and false.
	IsMature bool `json:"is_mature"`
}

type EventChannelUnban struct {
	// The user id for the user who was unbanned on the specified channel.
	UserID string `json:"user_id"`
	// The user login for the user who was unbanned on the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user who was unbanned on the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the issuer of the unban.
	ModeratorUserID string `json:"moderator_user_id"`
	// The user login of the issuer of the unban.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user name of the issuer of the unban.
	ModeratorUserName string `json:"moderator_user_name"`
}

type EventChannelFollow struct {
	// The user ID for the user now following the specified channel.
	UserID string `json:"user_id"`
	// The user login for the user now following the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user now following the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// RFC3339 timestamp of when the follow occurred.
	FollowedAt string `json:"followed_at"`
}

type EventChannelRaid struct {
	// The broadcaster ID that created the raid.
	FromBroadcasterUserID string `json:"from_broadcaster_user_id"`
	// The broadcaster login that created the raid.
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	// The broadcaster display name that created the raid.
	FromBroadcasterUserName string `json:"from_broadcaster_user_name"`
	// The broadcaster ID that received the raid.
	ToBroadcasterUserID string `json:"to_broadcaster_user_id"`
	// The broadcaster login that received the raid.
	ToBroadcasterUserLogin string `json:"to_broadcaster_user_login"`
	// The broadcaster display name that received the raid.
	ToBroadcasterUserName string `json:"to_broadcaster_user_name"`
	// The number of viewers in the raid.
	Viewers int `json:"viewers"`
}

type EventChannelModeratorAdd struct {
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the new moderator.
	UserID string `json:"user_id"`
	// The user login of the new moderator.
	UserLogin string `json:"user_login"`
	// The display name of the new moderator.
	UserName string `json:"user_name"`
}

type EventChannelModeratorRemove struct {
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the removed moderator.
	UserID string `json:"user_id"`
	// The user login of the removed moderator.
	UserLogin string `json:"user_login"`
	// The display name of the removed moderator.
	UserName string `json:"user_name"`
}

type EventChannelPollBegin struct {
	// ID of the poll.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Question displayed for the poll.
	Title string `json:"title"`
	// An array of choices for the poll.
	Choices []PollChoice `json:"choices"`
	// The Bits voting settings for the poll.
	BitsVoting BitsVoting `json:"bits_voting"`
	// The Channel Points voting settings for the poll.
	ChannelPointsVoting ChannelPointsVoting `json:"channel_points_voting"`
	// The time the poll started.
	StartedAt string `json:"started_at"`
	// The time the poll will end.
	EndsAt string `json:"ends_at"`
}

type EventChannelPollProgress struct {
	// ID of the poll.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Question displayed for the poll.
	Title string `json:"title"`
	// An array of choices for the poll. Includes vote counts.
	Choices []PollChoice `json:"choices"`
	// The Bits voting settings for the poll.
	BitsVoting BitsVoting `json:"bits_voting"`
	// The Channel Points voting settings for the poll.
	ChannelPointsVoting ChannelPointsVoting `json:"channel_points_voting"`
	// The time the poll started.
	StartedAt string `json:"started_at"`
	// The time the poll will end.
	EndsAt string `json:"ends_at"`
}

type EventChannelPollEnd struct {
	// ID of the poll.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Question displayed for the poll.
	Title string `json:"title"`
	// An array of choices for the poll. Includes vote counts.
	Choices []PollChoice `json:"choices"`
	// The Bits voting settings for the poll.
	BitsVoting BitsVoting `json:"bits_voting"`
	// The Channel Points voting settings for the poll.
	ChannelPointsVoting ChannelPointsVoting `json:"channel_points_voting"`
	// The status of the poll. Valid values are completed, archived, and terminated.
	Status string `json:"status"`
	// The time the poll started.
	StartedAt string `json:"started_at"`
	// The time the poll ended.
	EndedAt string `json:"ended_at"`
}

type EventChannelPointsRewardAdd struct {
	// The reward identifier.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Is the reward currently enabled. If false, the reward won’t show up to viewers.
	IsEnabled bool `json:"is_enabled"`
	// Is the reward currently paused. If true, viewers can’t redeem.
	IsPaused bool `json:"is_paused"`
	// Is the reward currently in stock. If false, viewers can’t redeem.
	IsInStock bool `json:"is_in_stock"`
	// The reward title.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
	// Does the viewer need to enter information when redeeming the reward.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// Should redemptions be set to fulfilled status immediately when redeemed and skip the request queue instead of the normal unfulfilled status.
	ShouldRedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// Whether a maximum per stream is enabled and what the maximum is.
	MaxPerStream MaxPerStream `json:"max_per_stream"`
	// Whether a maximum per user per stream is enabled and what the maximum is.
	MaxPerUserPerStream MaxPerUserPerStream `json:"max_per_user_per_stream"`
	// Custom background color for the reward. Format: Hex with # prefix. Example: #FA1ED2.
	BackgroundColor string `json:"background_color"`
	// Set of custom images of 1x, 2x and 4x sizes for the reward. Can be null if no images have been uploaded.
	Image Image `json:"image"`
	// Set of default images of 1x, 2x and 4x sizes for the reward.
	DefaultImage Image `json:"default_image"`
	// Whether a cooldown is enabled and what the cooldown is in seconds.
	GlobalCooldown GlobalCooldown `json:"global_cooldown"`
	// Timestamp of the cooldown expiration. null if the reward isn’t on cooldown.
	CooldownExpiresAt string `json:"cooldown_expires_at"`
	// The number of redemptions redeemed during the current live stream. Counts against the max_per_stream limit.
	// null if the broadcasters stream isn’t live or max_per_stream isn’t enabled.
	RedemptionsRedeemedCurrentStream int `json:"redemptions_redeemed_current_stream"`
}

type EventChannelPointsRewardUpdate struct {
	// The reward identifier.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Is the reward currently enabled. If false, the reward won’t show up to viewers.
	IsEnabled bool `json:"is_enabled"`
	// Is the reward currently paused. If true, viewers can’t redeem.
	IsPaused bool `json:"is_paused"`
	// Is the reward currently in stock. If false, viewers can’t redeem.
	IsInStock bool `json:"is_in_stock"`
	// The reward title.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
	// Does the viewer need to enter information when redeeming the reward.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// Should redemptions be set to fulfilled status immediately when redeemed and skip the request queue instead of the normal unfulfilled status.
	ShouldRedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// Whether a maximum per stream is enabled and what the maximum is.
	MaxPerStream MaxPerStream `json:"max_per_stream"`
	// Whether a maximum per user per stream is enabled and what the maximum is.
	MaxPerUserPerStream MaxPerUserPerStream `json:"max_per_user_per_stream"`
	// Custom background color for the reward. Format: Hex with # prefix. Example: #FA1ED2.
	BackgroundColor string `json:"background_color"`
	// Set of custom images of 1x, 2x and 4x sizes for the reward. Can be null if no images have been uploaded.
	Image Image `json:"image"`
	// Set of default images of 1x, 2x and 4x sizes for the reward.
	DefaultImage Image `json:"default_image"`
	// Whether a cooldown is enabled and what the cooldown is in seconds.
	GlobalCooldown GlobalCooldown `json:"global_cooldown"`
	// Timestamp of the cooldown expiration. null if the reward isn’t on cooldown.
	CooldownExpiresAt string `json:"cooldown_expires_at"`
	// The number of redemptions redeemed during the current live stream. Counts against the max_per_stream limit.
	// null if the broadcasters stream isn’t live or max_per_stream isn’t enabled.
	RedemptionsRedeemedCurrentStream int `json:"redemptions_redeemed_current_stream"`
}

type EventChannelPointsRewardRemove struct {
	// The reward identifier.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Is the reward currently enabled. If false, the reward won’t show up to viewers.
	IsEnabled bool `json:"is_enabled"`
	// Is the reward currently paused. If true, viewers can’t redeem.
	IsPaused bool `json:"is_paused"`
	// Is the reward currently in stock. If false, viewers can’t redeem.
	IsInStock bool `json:"is_in_stock"`
	// The reward title.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
	// Does the viewer need to enter information when redeeming the reward.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// Should redemptions be set to fulfilled status immediately when redeemed and skip the request queue instead of the normal unfulfilled status.
	ShouldRedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// Whether a maximum per stream is enabled and what the maximum is.
	MaxPerStream MaxPerStream `json:"max_per_stream"`
	// Whether a maximum per user per stream is enabled and what the maximum is.
	MaxPerUserPerStream MaxPerUserPerStream `json:"max_per_user_per_stream"`
	// Custom background color for the reward. Format: Hex with # prefix. Example: #FA1ED2.
	BackgroundColor string `json:"background_color"`
	// Set of custom images of 1x, 2x and 4x sizes for the reward. Can be null if no images have been uploaded.
	Image Image `json:"image"`
	// Set of default images of 1x, 2x and 4x sizes for the reward.
	DefaultImage Image `json:"default_image"`
	// Whether a cooldown is enabled and what the cooldown is in seconds.
	GlobalCooldown GlobalCooldown `json:"global_cooldown"`
	// Timestamp of the cooldown expiration. null if the reward isn’t on cooldown.
	CooldownExpiresAt string `json:"cooldown_expires_at"`
	// The number of redemptions redeemed during the current live stream. Counts against the max_per_stream limit.
	// null if the broadcasters stream isn’t live or max_per_stream isn’t enabled.
	RedemptionsRedeemedCurrentStream int `json:"redemptions_redeemed_current_stream"`
}

type EventChannelPointsRewardRedemptionAdd struct {
	// The redemption identifier.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of the user that redeemed the reward.
	UserID string `json:"user_id"`
	// Login of the user that redeemed the reward.
	UserLogin string `json:"user_login"`
	// Display name of the user that redeemed the reward.
	UserName string `json:"user_name"`
	// The user input provided. Empty string if not provided.
	UserInput string `json:"user_input"`
	// Defaults to unfulfilled. Possible values are unknown, unfulfilled, fulfilled, and canceled.
	Status string `json:"status"`
	// Basic information about the reward that was redeemed, at the time it was redeemed.
	Reward Reward `json:"reward"`
	// RFC3339 timestamp of when the reward was redeemed.
	RedeemedAt string `json:"redeemed_at"`
}

type EventChannelPointsRewardRedemptionUpdate struct {
	// The redemption identifier.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of the user that redeemed the reward.
	UserID string `json:"user_id"`
	// Login of the user that redeemed the reward.
	UserLogin string `json:"user_login"`
	// Display name of the user that redeemed the reward.
	UserName string `json:"user_name"`
	// The user input provided. Empty string if not provided.
	UserInput string `json:"user_input"`
	// Will be fulfilled or canceled. Possible values are unknown, unfulfilled, fulfilled, and canceled.
	Status string `json:"status"`
	// Basic information about the reward that was redeemed, at the time it was redeemed.
	Reward Reward `json:"reward"`
	// RFC3339 timestamp of when the reward was redeemed.
	RedeemedAt string `json:"redeemed_at"`
}

type EventChannelPredictionBegin struct {
	// Channel Points Prediction ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the Channel Points Prediction.
	Title string `json:"title"`
	// An array of outcomes for the Channel Points Prediction.
	Outcomes []PredictionOutcome `json:"outcomes"`
	// The time the Channel Points Prediction started.
	StartedAt string `json:"started_at"`
	// The time the Channel Points Prediction will automatically lock.
	LocksAt string `json:"locks_at"`
}

type EventChannelPredictionProgress struct {
	// Channel Points Prediction ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the Channel Points Prediction.
	Title string `json:"title"`
	// An array of outcomes for the Channel Points Prediction. Includes top_predictors.
	Outcomes []PredictionOutcome `json:"outcomes"`
	// The time the Channel Points Prediction started.
	StartedAt string `json:"started_at"`
	// The time the Channel Points Prediction will automatically lock.
	LocksAt string `json:"locks_at"`
}

type EventChannelPredictionLock struct {
	// Channel Points Prediction ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the Channel Points Prediction.
	Title string `json:"title"`
	// An array of outcomes for the Channel Points Prediction. Includes top_predictors.
	Outcomes []PredictionOutcome `json:"outcomes"`
	// The time the Channel Points Prediction started.
	StartedAt string `json:"started_at"`
	// The time the Channel Points Prediction was locked.
	LockedAt string `json:"locked_at"`
}

type EventChannelPredictionEnd struct {
	// Channel Points Prediction ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the Channel Points Prediction.
	Title string `json:"title"`
	// ID of the winning outcome.
	WinningOutcomeID string `json:"winning_outcome_id"`
	// An array of outcomes for the Channel Points Prediction. Includes top_predictors.
	Outcomes []PredictionOutcome `json:"outcomes"`
	// The status of the Channel Points Prediction. Valid values are resolved and canceled.
	Status string `json:"status"`
	// The time the Channel Points Prediction started.
	StartedAt string `json:"started_at"`
	// The time the Channel Points Prediction ended.
	EndedAt string `json:"ended_at"`
}

type EventChannelSubscriptionEnd struct {
	// The user ID for the user whose subscription ended.
	UserID string `json:"user_id"`
	// The user login for the user whose subscription ended.
	UserLogin string `json:"user_login"`
	// The user display name for the user whose subscription ended.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The tier of the subscription that ended. Valid values are 1000, 2000, and 3000.
	Tier string `json:"tier"`
	// Whether the subscription was a gift.
	IsGift bool `json:"is_gift"`
}

type EventChannelSubscriptionGift struct {
	// The user ID of the user who sent the subscription gift. Set to null if it was an anonymous subscription gift.
	UserID string `json:"user_id"`
	// The user login of the user who sent the gift. Set to null if it was an anonymous subscription gift.
	UserLogin string `json:"user_login"`
	// The user display name of the user who sent the gift. Set to null if it was an anonymous subscription gift.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The number of subscriptions in the subscription gift.
	Total int `json:"total"`
	// The tier of subscriptions in the subscription gift.
	Tier string `json:"tier"`
	// The number of subscriptions gifted by this user in the channel.
	// This value is null for anonymous gifts or if the gifter has opted out of sharing this information.
	CumulativeTotal int `json:"cumulative_total"`
	// Whether the subscription gift was anonymous.
	IsAnonymous bool `json:"is_anonymous"`
}

type EventChannelSubscriptionMessage struct {
	// The user ID of the user who sent a resubscription chat message.
	UserID string `json:"user_id"`
	// The user login of the user who sent a resubscription chat message.
	UserLogin string `json:"user_login"`
	// The user display name of the user who a resubscription chat message.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The tier of the user’s subscription.
	Tier string `json:"tier"`
	// An object that contains the resubscription message and emote information needed to recreate the message.
	Message Message `json:"message"`
	// The total number of months the user has been subscribed to the channel.
	CumulativeTotal int `json:"cumulative_total"`
	// The number of consecutive months the user’s current subscription has been active.
	// This value is null if the user has opted out of sharing this information.
	StreakMonths int `json:"streak_months"`
	// The month duration of the subscription.
	DurationMonths int `json:"duration_months"`
}

type EventDropEntitlementGrant struct {
	// Individual event ID, as assigned by EventSub. Use this for de-duplicating messages.
	ID string `json:"id"`
	// Entitlement object.
	Data []struct {
		// The ID of the organization that owns the game that has Drops enabled.
		OrganizationID string `json:"organization_id"`
		// Twitch category ID of the game that was being played when this benefit was entitled.
		CategoryID string `json:"category_id"`
		// The category name.
		CategoryName string `json:"category_name"`
		// The campaign this entitlement is associated with.
		CampaignID string `json:"campaign_id"`
		// Twitch user ID of the user who was granted the entitlement.
		UserID string `json:"user_id"`
		// The user display name of the user who was granted the entitlement.
		UserName string `json:"user_name"`
		// The user login of the user who was granted the entitlement.
		UserLogin string `json:"user_login"`
		// Unique identifier of the entitlement. Use this to de-duplicate entitlements.
		EntitlementID string `json:"entitlement_id"`
		// Identifier of the Benefit.
		BenefitID string `json:"benefit_id"`
		// UTC timestamp in ISO format when this entitlement was granted on Twitch.
		CreatedAt string `json:"created_at"`
	} `json:"data"`
}

type EventBitsTransactionCreate struct {
	// Client ID of the extension.
	ExtensionClientID string `json:"extension_client_id"`
	// Transaction ID.
	ID string `json:"id"`
	// The transaction’s broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The transaction’s broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The transaction’s broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The transaction’s user ID.
	UserID string `json:"user_id"`
	// The transaction’s user login.
	UserLogin string `json:"user_login"`
	// The transaction’s user display name.
	UserName string `json:"user_name"`
	// Additional extension product information.
	Product Product `json:"product"`
}

type EventGoals struct {
	// An ID that identifies this event.
	ID string `json:"id"`
	// An ID that uniquely identifies the broadcaster.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster’s display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster’s user handle.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The type of goal. Possible values are: followerssubscriptions
	Type string `json:"type"`
	// A description of the goal, if specified. The description may contain a maximum of 40 characters.
	Description string `json:"description"`
	// A Boolean value that indicates whether the broadcaster achieved their goal.
	// Is true if the goal was achieved; otherwise, false. Only the channel.goal.end event includes this field.
	IsAchieved bool `json:"is_achieved"`
	// The current value.If the goal is to increase followers, this field is set to the current number of followers.
	// This number increases with new followers and decreases if users unfollow the channel.
	// For subscriptions, current_amount is increased and decreased by the points value associated with the subscription tier.
	// For example, if a tier-two subscription is worth 2 points, current_amount is increased or decreased by 2, not 1.
	CurrentAmount int `json:"current_amount"`
	// The goal’s target value. For example, if the broadcaster has 200 followers before creating the goal,
	// and their goal is to double that number, this field is set to 400.
	TargetAmount int `json:"target_amount"`
	// The UTC timestamp in RFC 3339 format, which indicates when the broadcaster created the goal.
	StartedAt string `json:"started_at"`
	// The UTC timestamp in RFC 3339 format, which indicates when the broadcaster ended the goal.
	// Only the channel.goal.end event includes this field.
	EndedAt string `json:"ended_at"`
}

type EventHypeTrainBegin struct {
	// The Hype Train ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Total points contributed to the Hype Train.
	Total int `json:"total"`
	// The number of points contributed to the Hype Train at the current level.
	Progress int `json:"progress"`
	// The number of points required to reach the next level.
	Goal int `json:"goal"`
	// The contributors with the most points contributed.
	TopContributions []TopContributor `json:"top_contributions"`
	// The most recent contribution.
	LastContribution LastContribution `json:"last_contribution"`
	// The time when the Hype Train started.
	StartedAt string `json:"started_at"`
	// The time when the Hype Train expires. The expiration is extended when the Hype Train reaches a new level.
	ExpiresAt string `json:"expires_at"`
}

type EventHypeTrainProgress struct {
	// The Hype Train ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The current level of the Hype Train.
	Level int `json:"level"`
	// Total points contributed to the Hype Train.
	Total int `json:"total"`
	// The number of points contributed to the Hype Train at the current level.
	Progress int `json:"progress"`
	// The number of points required to reach the next level.
	Goal int `json:"goal"`
	// The contributors with the most points contributed.
	TopContributions []TopContributor `json:"top_contributions"`
	// The most recent contribution.
	LastContribution LastContribution `json:"last_contribution"`
	// The time when the Hype Train started.
	StartedAt string `json:"started_at"`
	// The time when the Hype Train expires. The expiration is extended when the Hype Train reaches a new level.
	ExpiresAt string `json:"expires_at"`
}

type EventHypeTrainEnd struct {
	// The Hype Train ID.
	ID string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The current level of the Hype Train.
	Level int `json:"level"`
	// Total points contributed to the Hype Train.
	Total int `json:"total"`
	// The contributors with the most points contributed.
	TopContributions []TopContributor `json:"top_contributions"`
	// The time when the Hype Train started.
	StartedAt string `json:"started_at"`
	// The time when the Hype Train ended.
	EndedAt string `json:"ended_at"`
	// The time when the Hype Train cooldown ends so that the next Hype Train can start.
	CooldownEndsAt string `json:"cooldown_ends_at"`
}

type EventStreamOnline struct {
	// The id of the stream.
	ID string `json:"id"`
	// The broadcaster’s user id.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The stream type. Valid values are: live, playlist, watch_party, premiere, rerun.
	Type string `json:"type"`
	// The timestamp at which the stream went online at.
	StartedAt string `json:"started_at"`
}

type EventStreamOffline struct {
	// The broadcaster’s user id.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
}

type EventUserAuthorizationGrant struct {
	// The client_id of the application that was granted user access.
	ClientID string `json:"client_id"`
	// The user id for the user who has granted authorization for your client id.
	UserID string `json:"user_id"`
	// The user login for the user who has granted authorization for your client id.
	UserLogin string `json:"user_login"`
	// The user display name for the user who has granted authorization for your client id.
	UserName string `json:"user_name"`
}

type EventUserAuthorizationRevoke struct {
	// The client_id of the application with revoked user access.
	ClientID string `json:"client_id"`
	// The user id for the user who has revoked authorization for your client id.
	UserID string `json:"user_id"`
	// The user login for the user who has revoked authorization for your client id. This is null if the user no longer exists.
	UserLogin string `json:"user_login"`
	// The user display name for the user who has revoked authorization for your client id. This is null if the user no longer exists.
	UserName string `json:"user_name"`
}

type EventUserUpdate struct {
	// The user’s user id.
	UserID string `json:"user_id"`
	// The user’s user login.
	UserLogin string `json:"user_login"`
	// The user’s user display name.
	UserName string `json:"user_name"`
	// The user’s email. Only included if you have the user:read:email scope for the user.
	Email string `json:"email"`
	// The user’s description.
	Description string `json:"description"`
}
