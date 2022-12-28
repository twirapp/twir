package modules

type YoutubeUserSettings struct {
	MaxRequests   int   ` json:"maxRequests"`
	MinWatchTime  int64 ` json:"minWatchTime"`
	MinMessages   int   `json:"minMessages"`
	MinFollowTime int   ` json:"minFollowTime"`
}

type YotubeSongSettings struct {
	MaxLength          int      `validate:"lte=86400"          json:"maxLength"`
	MinViews           int      `validate:"lte=10000000000000" json:"minViews"`
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
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YoutubeBlacklistSettings struct {
	Users        []YoutubeBlacklistSettingsUsers    `validate:"required,dive"         json:"users"`
	Songs        []YoutubeBlacklistSettingsSongs    `validate:"required,dive"         json:"songs"`
	Channels     []YoutubeBlacklistSettingsChannels `validate:"required,dive"         json:"channels"`
	ArtistsNames []string                           `validate:"required,dive,max=300" json:"artistsNames"`
}

type YoutubeSettings struct {
	Enabled               *bool                    `validate:"required" json:"enabled"`
	AcceptOnlyWhenOnline  *bool                    `validate:"required"                   json:"acceptOnlyWhenOnline"`
	MaxRequests           int                      `validate:"lte=500" json:"maxRequests"`
	ChannelPointsRewardId string                   `validate:"max=100" json:"channelPointsRewardId"`
	User                  YoutubeUserSettings      `validate:"required"    json:"user"`
	Song                  YotubeSongSettings       `validate:"required"    json:"song"`
	BlackList             YoutubeBlacklistSettings `validate:"required"    json:"blackList"`
}

type SearchResult struct {
	ID        string `     json:"id"`
	Title     string ` json:"title"`
	ThumbNail string `json:"thumbNail"`
}

type YouTube struct {
	POST                    YoutubeSettings
	GET                     YoutubeSettings
	POST_BLACKLIST_SONGS    YoutubeBlacklistSettingsSongs
	POST_BLACKLIST_CHANNELS YoutubeBlacklistSettingsChannels
	POST_BLACKLIST_USERS    YoutubeBlacklistSettingsUsers
	SEARCH                  []SearchResult
}
