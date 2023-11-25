package model

type ChannelModulesSettingsBeRightBack struct {
	Text string `json:"text"`

	Late    ChannelModulesSettingsBeRightBackLate `json:"late"`
	Opacity int32                                 `json:"opacity"`

	BackgroundColor string `json:"backgroundColor"`

	FontSize  int32  `json:"fontSize"`
	FontColor string `json:"color"`
}

type ChannelModulesSettingsBeRightBackLate struct {
	Enabled        bool   `json:"enabled"`
	Text           string `json:"text"`
	DisplayBrbTime bool   `json:"displayBrbTime"`
}
