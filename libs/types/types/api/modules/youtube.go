package modules

type YoutubeUserSettings struct {
	MaxRequests   *int   `json:"maxRequests"`
	MinWatchTime  *int64 `json:"minWatchTime"`
	MinMessages   *int   `json:"minMessages"`
	MinFollowTime *int   `json:"minFollowTime"`
}

type YotubeSongSettings struct {
	MaxLength          *int     `validate:"lte=86400"          json:"maxLength"`
	MinViews           *int     `validate:"lte=10000000000000" json:"minViews"`
	AcceptedCategories []string `validate:"dive,max=300"       json:"acceptedCategories"`
}

type YoutubeBlacklistSettingsUsers struct {
	UserID   string `json:"userId"   validate:"max=50"`
	UserName string `json:"userName" validate:"required"`
}

type YoutubeBlacklistSettingsSongs struct {
	ID        string `validate:"required,min=1,max=300" json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YoutubeBlacklistSettingsChannels struct {
	ID        string `validate:"required,min=1"          json:"id"`
	Title     string `validate:"required,min=1,max=300"  json:"title"`
	ThumbNail string `validate:"required,min=1=,max=300" json:"thumbNail"`
}

type YoutubeBlacklistSettings struct {
	Users        []YoutubeBlacklistSettingsUsers    `validate:"dive"         json:"users"`
	Songs        []YoutubeBlacklistSettingsSongs    `validate:"dive"         json:"songs"`
	Channels     []YoutubeBlacklistSettingsChannels `validate:"dive"         json:"channels"`
	ArtistsNames []string                           `validate:"dive,max=300" json:"artistsNames"`
}

type YoutubeSettings struct {
	MaxRequests             *int                      `validate:"lte=500" json:"maxRequests"`
	AcceptOnlyWhenOnline    *bool                     `                   json:"acceptOnlyWhenOnline"`
	ChannelPointsRewardName *string                   `validate:"max=100" json:"channelPointsRewardName"`
	User                    *YoutubeUserSettings      `validate:"dive"    json:"user"`
	Song                    *YotubeSongSettings       `validate:"dive"    json:"song"`
	BlackList               *YoutubeBlacklistSettings `validate:"dive"    json:"blacklist"`
}

type YouTube struct {
	POST                    YoutubeSettings
	GET                     YoutubeSettings
	POST_BLACKLIST_SONGS    YoutubeBlacklistSettingsSongs
	POST_BLACKLIST_CHANNELS YoutubeBlacklistSettingsChannels
	POST_BLACKLIST_USERS    YoutubeBlacklistSettingsUsers
}
