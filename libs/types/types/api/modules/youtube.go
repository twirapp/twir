package modules

type YouTubeUserSettings struct {
	MaxRequests   int   `json:"maxRequests"`
	MinWatchTime  int64 `json:"minWatchTime"`
	MinMessages   int   `json:"minMessages"`
	MinFollowTime int   `json:"minFollowTime"`
}

type YouTubeSongSettings struct {
	MaxLength          int      `validate:"lte=86400"          json:"maxLength"`
	MinViews           int      `validate:"lte=10000000000000" json:"minViews"`
	AcceptedCategories []string `validate:"dive,max=300"       json:"acceptedCategories"`
}

type YouTubeDenySettingsUsers struct {
	UserID   string `json:"userId"   validate:"max=50"`
	UserName string `json:"userName" validate:"required"`
}

type YouTubeDenySettingsSongs struct {
	ID        string `validate:"required,min=1,max=300" json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YouTubeDenySettingsChannels struct {
	ID        string `validate:"required,min=1"         json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YouTubeDenyList struct {
	Users        []YouTubeDenySettingsUsers    `validate:"required,dive"         json:"users"`
	Songs        []YouTubeDenySettingsSongs    `validate:"required,dive"         json:"songs"`
	Channels     []YouTubeDenySettingsChannels `validate:"required,dive"         json:"channels"`
	ArtistsNames []string                      `validate:"required,dive,max=300" json:"artistsNames"`
}

type YouTubeSettings struct {
	Enabled               *bool               `validate:"required" json:"enabled"`
	AcceptOnlyWhenOnline  *bool               `validate:"required" json:"acceptOnlyWhenOnline"`
	MaxRequests           int                 `validate:"lte=500"  json:"maxRequests"`
	ChannelPointsRewardId string              `validate:"max=100"  json:"channelPointsRewardId"`
	AnnouncePlay          *bool               `validate:"required" json:"announcePlay"`
	User                  YouTubeUserSettings `validate:"required" json:"user"`
	Song                  YouTubeSongSettings `validate:"required" json:"song"`
	DenyList              YouTubeDenyList     `validate:"required" json:"denyList"`
}

type SearchResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ThumbNail string `json:"thumbNail"`
}

type YouTube struct {
	POST   YouTubeSettings
	GET    YouTubeSettings
	SEARCH []SearchResult
}
