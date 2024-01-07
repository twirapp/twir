package model

type ChannelModulesSettingsBeRightBack struct {
	Text string `json:"text"`

	Late ChannelModulesSettingsBeRightBackLate `json:"late"`

	BackgroundColor string `json:"backgroundColor"`

	FontSize   int32  `json:"fontSize"`
	FontColor  string `json:"fontColor"`
	FontFamily string `json:"fontFamily"`
}

type ChannelModulesSettingsBeRightBackLate struct {
	Enabled        bool   `json:"enabled"`
	Text           string `json:"text"`
	DisplayBrbTime bool   `json:"displayBrbTime"`
}
