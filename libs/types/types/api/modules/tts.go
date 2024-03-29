package modules

type TTSSettings struct {
	Enabled                            *bool    `validate:"required" json:"enabled"`
	Rate                               int      `validate:"gte=0,lte=100" json:"rate"`
	Volume                             int      `validate:"gte=0,lte=100" json:"volume"`
	Pitch                              int      `validate:"gte=0,lte=100" json:"pitch"`
	Voice                              string   `validate:"required" json:"voice"`
	AllowUsersChooseVoiceInMainCommand bool     `json:"allow_users_choose_voice_in_main_command"`
	MaxSymbols                         int      `json:"max_symbols"`
	DisallowedVoices                   []string `json:"disallowed_voices"`
	DoNotReadEmoji                     bool     `json:"do_not_read_emoji"`
	DoNotReadTwitchEmotes              bool     `json:"do_not_read_twitch_emotes"`
	DoNotReadLinks                     bool     `json:"do_not_read_links"`
	ReadChatMessages                   bool     `json:"read_chat_messages"`
	ReadChatMessagesNicknames          bool     `json:"read_chat_messages_nicknames"`
}

type TTS struct {
	GET  TTSSettings
	POST TTSSettings
}
