package entity

type TwitchCustomReward struct {
	BroadcasterID                     string                                        `json:"broadcaster_id"`
	BroadcasterLogin                  string                                        `json:"broadcaster_login"`
	BroadcasterName                   string                                        `json:"broadcaster_name"`
	ID                                string                                        `json:"id"`
	Title                             string                                        `json:"title"`
	Prompt                            string                                        `json:"prompt"`
	Cost                              int                                           `json:"cost"`
	Image                             TwitchCustomRewardImage                       `json:"image"`
	BackgroundColor                   string                                        `json:"background_color"`
	DefaultImage                      TwitchCustomRewardImage                       `json:"default_image"`
	IsEnabled                         bool                                          `json:"is_enabled"`
	IsUserInputRequired               bool                                          `json:"is_user_input_required"`
	MaxPerStreamSetting               TwitchCustomRewardMaxPerStreamSettings        `json:"max_per_stream_setting"`
	MaxPerUserPerStreamSetting        TwitchCustomRewardMaxPerUserPerStreamSettings `json:"max_per_user_per_stream_setting"`
	GlobalCooldownSetting             TwitchCustomRewardGlobalCooldownSettings      `json:"global_cooldown_setting"`
	IsPaused                          bool                                          `json:"is_paused"`
	IsInStock                         bool                                          `json:"is_in_stock"`
	ShouldRedemptionsSkipRequestQueue bool                                          `json:"should_redemptions_skip_request_queue"`
	RedemptionsRedeemedCurrentStream  int                                           `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 string                                        `json:"cooldown_expires_at"`
}

type TwitchCustomRewardImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

type TwitchCustomRewardMaxPerUserPerStreamSettings struct {
	IsEnabled           bool `json:"is_enabled"`
	MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
}

type TwitchCustomRewardMaxPerStreamSettings struct {
	IsEnabled    bool `json:"is_enabled"`
	MaxPerStream int  `json:"max_per_stream"`
}

type TwitchCustomRewardGlobalCooldownSettings struct {
	IsEnabled             bool `json:"is_enabled"`
	GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
}
