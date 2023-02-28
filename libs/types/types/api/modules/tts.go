package modules

type TTSSettings struct {
	Enabled                            *bool  `validate:"required" json:"enabled"`
	Rate                               int    `validate:"gte=0,lte=100" json:"rate"`
	Volume                             int    `validate:"gte=0,lte=100" json:"volume"`
	Pitch                              int    `validate:"gte=0,lte=100" json:"pitch"`
	Voice                              string `validate:"required" json:"voice"`
	AllowUsersChooseVoiceInMainCommand bool   `validate:"required" json:"allow_users_choose_voice_in_main_command"`
	MaxSymbols                         int    `json:"max_symbols"`
}

type TTS struct {
	GET  TTSSettings
	POST TTSSettings
}
