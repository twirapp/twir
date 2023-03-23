package eventsub_bindings

type PollChoice struct {
	// ID for the choice.
	ID string `json:"id"`
	// Text displayed for the choice.
	Title string `json:"title"`
	// Number of votes received via Bits.
	BitsVotes int `json:"bits_votes"`
	// Number of votes received via Channel Points.
	ChannelPointsVotes int `json:"channel_points_votes"`
	// Total number of votes received for the choice across all methods of voting.
	Votes int `json:"votes"`
}

type BitsVoting struct {
	// Indicates if Bits can be used for voting.
	IsEnabled bool `json:"is_enabled"`
	// Number of Bits required to vote once with Bits.
	AmountPerVote int `json:"amount_per_vote"`
}

type ChannelPointsVoting struct {
	// Indicates if Channel Points can be used for voting.
	IsEnabled bool `json:"is_enabled"`
	// Number of Channel Points required to vote once with Channel Points.
	AmountPerVote int `json:"amount_per_vote"`
}

type MaxPerStream struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The max per stream limit.
	Value int `json:"value"`
}

type MaxPerUserPerStream struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The max per user per stream limit.
	Value int `json:"value"`
}

type Image struct {
	// URL for the image at 1x size.
	Url1x string `json:"url_1x"`
	// URL for the image at 2x size.
	Url2x string `json:"url_2x"`
	// URL for the image at 4x size.
	Url4x string `json:"url_4x"`
}

type GlobalCooldown struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The cooldown in seconds.
	Seconds int `json:"seconds"`
}

type Reward struct {
	// The reward identifier.
	ID string `json:"id"`
	// The reward name.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
}

type PredictionOutcome struct {
	// The outcome ID.
	ID string `json:"id"`
	// The outcome title.
	Title string `json:"title"`
	// The color for the outcome. Valid values are pink and blue.
	Color string `json:"color"`
	// The number of users who used Channel Points on this outcome.
	Users int `json:"users"`
	// The total number of Channel Points used on this outcome.
	ChannelPoints int `json:"channel_points"`
	// An array of users who used the most Channel Points on this outcome.
	TopPredictors []TopPredictor `json:"top_predictors"`
}

type TopPredictor struct {
	// The ID of the user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// The number of Channel Points won.
	// This value is always null in the event payload for Prediction progress and Prediction lock.
	// This value is 0 if the outcome did not win or if the Prediction was canceled and Channel Points were refunded.
	ChannelPointsWon int `json:"channel_points_won"`
	// The number of Channel Points used to participate in the Prediction.
	ChannelPointsUsed int `json:"channel_points_used"`
}

type Message struct {
	// The text of the resubscription chat message.
	Text string `json:"text"`
	// An array that includes the emote ID and start and end positions for where the emote appears in the text.
	Emotes []Emote `json:"emotes"`
}

type Emote struct {
	// The index of where the Emote starts in the text.
	Begin int `json:"begin"`
	// The index of where the Emote ends in the text.
	End int `json:"end"`
	// The emote ID.
	ID string `json:"id"`
}

type Product struct {
	// Product name.
	Name string `json:"name"`
	// Bits involved in the transaction.
	Bits int `json:"bits"`
	// Unique identifier for the product acquired.
	Sku string `json:"sku"`
	// Flag indicating if the product is in development. If InDevelopment is true, bits will be 0.
	InDevelopment bool `json:"in_development"`
}

type TopContributor struct {
	// The ID of the user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// Type of contribution. Valid values include bits, subscription.
	Type string `json:"type"`
	// The total contributed.
	Total int `json:"total"`
}

type LastContribution struct {
	// The ID of the user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// Type of contribution. Valid values include bits, subscription.
	Type string `json:"type"`
	// The total contributed.
	Total int `json:"total"`
}
