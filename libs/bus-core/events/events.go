package events

const (
	FollowSubject                     = "events.follow"
	SubscribeSubject                  = "events.subscribe"
	SubGiftSubject                    = "events.sub_gift"
	ReSubscribeSubject                = "events.resubscribe"
	RedemptionCreatedSubject          = "events.redemption_created"
	CommandUsedSubject                = "events.command_used"
	FirstUserMessageSubject           = "events.first_user_message"
	RaidedSubject                     = "events.raided"
	TitleOrCategoryChangedSubject     = "events.title_or_category_changed"
	StreamOnlineSubject               = "events.stream_online"
	StreamOfflineSubject              = "events.stream_offline"
	ChatClearSubject                  = "events.chat_clear"
	DonateSubject                     = "events.donate"
	KeywordMatchedSubject             = "events.keyword_matched"
	GreetingSendedSubject             = "events.greeting_sended"
	PollBeginSubject                  = "events.poll_begin"
	PollProgressSubject               = "events.poll_progress"
	PollEndSubject                    = "events.poll_end"
	PredictionBeginSubject            = "events.prediction_begin"
	PredictionProgressSubject         = "events.prediction_progress"
	PredictionLockSubject             = "events.prediction_lock"
	PredictionEndSubject              = "events.prediction_end"
	StreamFirstUserJoinSubject        = "events.stream_first_user_join"
	ChannelBanSubject                 = "events.channel_ban"
	ChannelUnbanRequestCreateSubject  = "events.channel_unban_request_create"
	ChannelUnbanRequestResolveSubject = "events.channel_unban_request_resolve"
	ChannelMessageDeleteSubject       = "events.channel_message_delete"
	VipAddedSubject                   = "events.vip_added"
	VipRemovedSubject                 = "events.vip_removed"
	ModeratorAddedSubject             = "events.moderator_added"
	ModeratorRemovedSubject           = "events.moderator_removed"
	ChannelUnbanSubject               = "events.channel_unban"
)

type BaseInfo struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type FollowMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	UserID          string   `json:"user_id"`
}

type SubscribeMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	Level           string   `json:"level"`
	UserID          string   `json:"user_id"`
}

type SubGiftMessage struct {
	BaseInfo          BaseInfo `json:"base_info"`
	SenderUserName    string   `json:"sender_user_name"`
	SenderDisplayName string   `json:"sender_display_name"`
	TargetUserName    string   `json:"target_user_name"`
	TargetDisplayName string   `json:"target_display_name"`
	Level             string   `json:"level"`
	SenderUserID      string   `json:"sender_user_id"`
}

type ReSubscribeMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	Months          int64    `json:"months"`
	Streak          int64    `json:"streak"`
	IsPrime         bool     `json:"is_prime"`
	Message         string   `json:"message"`
	Level           string   `json:"level"`
	UserID          string   `json:"user_id"`
}

type RedemptionCreatedMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	ID              string   `json:"id"`
	RewardName      string   `json:"reward_name"`
	RewardCost      string   `json:"reward_cost"`
	Input           *string  `json:"input,omitempty"`
	UserID          string   `json:"user_id"`
}

type CommandUsedMessage struct {
	BaseInfo           BaseInfo `json:"base_info"`
	CommandID          string   `json:"command_id"`
	CommandName        string   `json:"command_name"`
	UserName           string   `json:"user_name"`
	UserDisplayName    string   `json:"user_display_name"`
	CommandInput       string   `json:"command_input"`
	UserID             string   `json:"user_id"`
	IsDefault          bool     `json:"is_default"`
	DefaultCommandName *string  `json:"default_command_name"`
	MessageID          string   `json:"message_id"`
}

type FirstUserMessageMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserID          string   `json:"user_id"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	MessageID       string   `json:"message_id"`
}

type RaidedMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	Viewers         int64    `json:"viewers"`
	UserID          string   `json:"user_id"`
}

type TitleOrCategoryChangedMessage struct {
	BaseInfo    BaseInfo `json:"base_info"`
	OldTitle    string   `json:"old_title"`
	NewTitle    string   `json:"new_title"`
	OldCategory string   `json:"old_category"`
	NewCategory string   `json:"new_category"`
}

type ChatClearMessage struct {
	BaseInfo BaseInfo `json:"base_info"`
}

type DonateMessage struct {
	BaseInfo BaseInfo `json:"base_info"`
	UserName string   `json:"user_name"`
	Amount   string   `json:"amount"`
	Currency string   `json:"currency"`
	Message  string   `json:"message"`
}

type KeywordMatchedMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	KeywordID       string   `json:"keyword_id"`
	KeywordName     string   `json:"keyword_name"`
	KeywordResponse string   `json:"keyword_response"`
	UserID          string   `json:"user_id"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
}

type GreetingSendedMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserID          string   `json:"user_id"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	GreetingText    string   `json:"greeting_text"`
}

type PollChoice struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	BitsVotes           uint64 `json:"bits_votes"`
	ChannelsPointsVotes uint64 `json:"channels_points_votes"`
	Votes               uint64 `json:"votes"`
}

type PollBitsVotes struct {
	Enabled       bool   `json:"enabled"`
	AmountPerVote uint64 `json:"amount_per_vote"`
}

type PollChannelPointsVotes struct {
	Enabled       bool   `json:"enabled"`
	AmountPerVote uint64 `json:"amount_per_vote"`
}

type PollInfo struct {
	Title                string                 `json:"title"`
	Choices              []PollChoice           `json:"choices"`
	BitsVoting           PollBitsVotes          `json:"bits_voting"`
	ChannelsPointsVoting PollChannelPointsVotes `json:"channels_points_voting"`
}

type PollBeginMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	Info            PollInfo `json:"info"`
}

type PollProgressMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	Info            PollInfo `json:"info"`
}

type PollEndMessage struct {
	BaseInfo        BaseInfo `json:"base_info"`
	UserName        string   `json:"user_name"`
	UserDisplayName string   `json:"user_display_name"`
	Info            PollInfo `json:"info"`
}

type PredictionTopPredictor struct {
	UserName        string  `json:"user_name"`
	UserDisplayName string  `json:"user_display_name"`
	UserID          string  `json:"user_id"`
	PointsUsed      uint64  `json:"points_used"`
	PointsWin       *uint64 `json:"points_win,omitempty"`
}

type PredictionOutcome struct {
	ID            string                   `json:"id"`
	Title         string                   `json:"title"`
	Color         string                   `json:"color"`
	Users         uint64                   `json:"users"`
	ChannelPoints uint64                   `json:"channel_points"`
	TopPredictors []PredictionTopPredictor `json:"top_predictors"`
}

type PredictionInfo struct {
	Title    string              `json:"title"`
	Outcomes []PredictionOutcome `json:"outcomes"`
}

type PredictionBeginMessage struct {
	BaseInfo        BaseInfo       `json:"base_info"`
	UserName        string         `json:"user_name"`
	UserDisplayName string         `json:"user_display_name"`
	Info            PredictionInfo `json:"info"`
}

type PredictionProgressMessage struct {
	BaseInfo        BaseInfo       `json:"base_info"`
	UserName        string         `json:"user_name"`
	UserDisplayName string         `json:"user_display_name"`
	Info            PredictionInfo `json:"info"`
}

type PredictionLockMessage struct {
	BaseInfo        BaseInfo       `json:"base_info"`
	UserName        string         `json:"user_name"`
	UserDisplayName string         `json:"user_display_name"`
	Info            PredictionInfo `json:"info"`
}

type PredictionEndMessage struct {
	BaseInfo         BaseInfo       `json:"base_info"`
	UserName         string         `json:"user_name"`
	UserDisplayName  string         `json:"user_display_name"`
	Info             PredictionInfo `json:"info"`
	WinningOutcomeID string         `json:"winning_outcome_id"`
}

type StreamFirstUserJoinMessage struct {
	BaseInfo  BaseInfo `json:"base_info"`
	UserID    string   `json:"user_id"`
	UserLogin string   `json:"user_login"`
}

type ChannelBanMessage struct {
	BaseInfo             BaseInfo `json:"base_info"`
	UserID               string   `json:"user_id"`
	UserName             string   `json:"user_name"`
	UserLogin            string   `json:"user_login"`
	BroadcasterUserName  string   `json:"broadcaster_user_name"`
	BroadcasterUserLogin string   `json:"broadcaster_user_login"`
	ModeratorUserID      string   `json:"moderator_id"`
	ModeratorUserName    string   `json:"moderator_user_name"`
	ModeratorUserLogin   string   `json:"moderator_user_login"`
	Reason               string   `json:"reason"`
	EndsAt               string   `json:"ends_at"`
	IsPermanent          bool     `json:"is_permanent"`
}

type ChannelUnbanMessage struct {
	BaseInfo             BaseInfo `json:"base_info"`
	UserID               string   `json:"user_id"`
	UserName             string   `json:"user_name"`
	UserLogin            string   `json:"user_login"`
	BroadcasterUserName  string   `json:"broadcaster_user_name"`
	BroadcasterUserLogin string   `json:"broadcaster_user_login"`
	ModeratorUserID      string   `json:"moderator_id"`
	ModeratorUserName    string   `json:"moderator_user_name"`
	ModeratorUserLogin   string   `json:"moderator_user_login"`
}

type ChannelUnbanRequestCreateMessage struct {
	BaseInfo             BaseInfo `json:"base_info"`
	UserName             string   `json:"user_name"`
	UserLogin            string   `json:"user_login"`
	BroadcasterUserName  string   `json:"broadcaster_user_name"`
	BroadcasterUserLogin string   `json:"broadcaster_user_login"`
	Text                 string   `json:"text"`
}

type ChannelUnbanRequestResolveMessage struct {
	BaseInfo             BaseInfo `json:"base_info"`
	UserID               string   `json:"user_id"`
	UserName             string   `json:"user_name"`
	UserLogin            string   `json:"user_login"`
	BroadcasterUserName  string   `json:"broadcaster_user_name"`
	BroadcasterUserLogin string   `json:"broadcaster_user_login"`
	ModeratorUserID      string   `json:"moderator_id"`
	ModeratorUserName    string   `json:"moderator_user_name"`
	ModeratorUserLogin   string   `json:"moderator_user_login"`
	Declined             bool     `json:"declined"`
	Reason               string   `json:"reason"`
}

type ChannelMessageDeleteMessage struct {
	BaseInfo             BaseInfo `json:"base_info"`
	UserId               string   `json:"user_id"`
	UserName             string   `json:"user_name"`
	UserLogin            string   `json:"user_login"`
	BroadcasterUserName  string   `json:"broadcaster_user_name"`
	BroadcasterUserLogin string   `json:"broadcaster_user_login"`
	MessageId            string   `json:"message_id"`
}

type VipAddedMessage struct {
	BaseInfo BaseInfo `json:"base_info"`
	UserName string   `json:"user_name"`
	UserID   string   `json:"user_id"`
}

type VipRemovedMessage struct {
	BaseInfo BaseInfo `json:"base_info"`
	UserName string   `json:"user_name"`
	UserID   string   `json:"user_id"`
}

type ModeratorAddedMessage struct {
	BaseInfo BaseInfo `json:"base_info"`
	UserName string   `json:"user_name"`
	UserID   string   `json:"user_id"`
}

type ModeratorRemovedMessage struct {
	BaseInfo BaseInfo `json:"base_info"`
	UserName string   `json:"user_name"`
	UserID   string   `json:"user_id"`
}
