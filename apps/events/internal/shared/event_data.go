package shared

type EventData struct {
	ChannelID string `json:"channelId"`

	//
	UserName        string `json:"userName,omitempty"`
	UserDisplayName string `json:"userDisplayName,omitempty"`
	UserID          string `json:"userId,omitempty"`
	//
	RaidViewers int64 `json:"raidViewers,omitempty"`
	//
	ResubMonths  int64  `json:"resubMonths,omitempty"`
	ResubStreak  int64  `json:"resubStreak,omitempty"`
	ResubMessage string `json:"resubMessage,omitempty"`
	SubLevel     string `json:"subLevel,omitempty"`
	//
	OldStreamTitle    string `json:"oldStreamTitle,omitempty"`
	NewStreamTitle    string `json:"newStreamTitle,omitempty"`
	OldStreamCategory string `json:"oldStreamCategory,omitempty"`
	NewStreamCategory string `json:"newStreamCategory,omitempty"`
	//
	StreamTitle    string `json:"streamTitle,omitempty"`
	StreamCategory string `json:"streamCategory,omitempty"`
	//
	RewardID    string  `json:"rewardId,omitempty"`
	RewardName  string  `json:"rewardName,omitempty"`
	RewardCost  string  `json:"rewardCost,omitempty"`
	RewardInput *string `json:"rewardInput,omitempty"`
	//
	CommandName  string `json:"commandName,omitempty"`
	CommandID    string `json:"commandId,omitempty"`
	CommandInput string `json:"commandInput,omitempty"`

	//
	TargetUserName        string `json:"targetUserName,omitempty"`
	TargetUserDisplayName string `json:"targetUserDisplayName,omitempty"`
	//
	DonateAmount   string `json:"donateAmount,omitempty"`
	DonateMessage  string `json:"donateMessage,omitempty"`
	DonateCurrency string `json:"donateCurrency,omitempty"`
	//

	//
	KeywordID       string `json:"keywordId,omitempty"`
	KeywordName     string `json:"keywordName,omitempty"`
	KeywordResponse string `json:"keywordResponse,omitempty"`
	//
	GreetingText string `json:"greetingText,omitempty"`
	//
	PollTitle                     string `json:"pollTitle,omitempty"`
	PollOptionsNames              string `json:"pollOptionsNames,omitempty"`
	PollTotalVotes                int    `json:"pollTotalVotes,omitempty"`
	PollWinnerTitle               string `json:"pollWinnerTitle,omitempty"`
	PollWinnerBitsVotes           int    `json:"pollWinnerBitsVotes,omitempty"`
	PollWinnerChannelsPointsVotes int    `json:"pollWinnerChannelsPointsVotes,omitempty"`
	PollWinnerTotalVotes          int    `json:"pollWinnerTotalVotes,omitempty"`
	//
	PredictionTitle              string                      `json:"predictionTitle,omitempty"`
	PredictionOptionsNames       string                      `json:"predictionOptionsNames,omitempty"`
	PredictionTotalChannelPoints int                         `json:"predictionTotalChannelPoints,omitempty"`
	PredictionWinner             *EventDataPredictionOutCome `json:"predictionWinner,omitempty"`
	//

	ModeratorName        string `json:"moderatorName,omitempty"`
	ModeratorDisplayName string `json:"moderatorDisplayName,omitempty"`

	// ban
	BanReason        string `json:"banReason,omitempty"`
	BanEndsInMinutes string `json:"banEndsInMinutes,omitempty"`

	// generic message field
	Message string `json:"message,omitempty"`

	//
	ChannelUnbanRequestResolveDeclined bool `json:"channelUnbanRequestResolveStatus,omitempty"`
}

type EventDataPredictionOutCome struct {
	Title       string `json:"title,omitempty"`
	TotalUsers  int    `json:"totalUsers,omitempty"`
	TotalPoints int    `json:"totalPoints,omitempty"`
	TopUsers    string `json:"topUsers,omitempty"`
}

// type DataFromPrevOperation struct {
// 	UnmodedUserName string `json:"unmodedUserName"`
// 	UnvipedUserName string `json:"unvipedUserName"`
// 	BannedUserName  string `json:"bannedUserName"`
// }
