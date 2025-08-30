package model

type ChannelModulesSettingsTTS struct {
	ID                                 string   `json:"id"`
	Type                               string   `json:"type"` // should be always "tts"
	ChannelID                          string   `json:"channelId"`
	UserID                             *string  `json:"userId"`
	Enabled                            *bool    `json:"enabled"`
	Rate                               int      `json:"rate"`
	Volume                             int      `json:"volume"`
	Pitch                              int      `json:"pitch"`
	Voice                              string   `json:"voice"`
	AllowUsersChooseVoiceInMainCommand bool     `json:"allow_users_choose_voice_in_main_command"`
	MaxSymbols                         int      `json:"max_symbols"`
	DisallowedVoices                   []string `json:"disallowed_voices"`
	DoNotReadEmoji                     bool     `json:"do_not_read_emoji"`
	DoNotReadTwitchEmotes              bool     `json:"do_not_read_twitch_emotes"`
	DoNotReadLinks                     bool     `json:"do_not_read_links"`
	ReadChatMessages                   bool     `json:"read_chat_messages"`
	ReadChatMessagesNicknames          bool     `json:"read_chat_messages_nicknames"`
}

var Nil = ChannelModulesSettingsTTS{}
