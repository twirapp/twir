package types

// export type YoutubeSettings = {
//   maxRequests?: number;
//   acceptOnlyWhenOnline?: boolean;
//   channelPointsRewardName?: string;
//   user?: {
//     maxRequests?: number;
//     minWatchTime?: number;
//     minMessages?: number;
//     minFollowTime?: number;
//   };
//   song?: {
//     maxLength?: number;
//     minViews?: number;
//     acceptedCategories?: string[];
//   };
//   blackList?: {
//     users?: Array<{ userId: string, userName: string }>;
//     songs?: Array<{ id: string, title: string, thumbnail: string }>;
//     channels?: Array<{ id: string, title: string, thumbnail: string }>;
//     artistsNames?: string[];
//   };
// };

type YoutubeUserSettings struct {
	MaxRequests   *int `json:"maxRequests"`
	MinWatchTime  *int `json:"minWatchTime"`
	MinMessages   *int `json:"minMessages"`
	MinFollowTime *int `json:"minFollowTime"`
}

type YotubeSongSettings struct {
	MaxLength          *int     `json:"maxLength"`
	MinViews           *int     `json:"minViews"`
	AcceptedCategories []string `json:"acceptedCategories"`
}

type YoutubeBlacklistSettingsUsers struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
}

type YoutubeBlacklistSettingsSongs struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ThumbNail string `json:"thumbNail"`
}

type YoutubeBlacklistSettingsChannels struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ThumbNail string `json:"thumbNail"`
}

type YoutubeBlacklistSettings struct {
	Users        []YoutubeBlacklistSettingsUsers    `json:"users"`
	Songs        []YoutubeBlacklistSettingsSongs    `json:"songs"`
	Channels     []YoutubeBlacklistSettingsChannels `json:"channels"`
	ArtistsNames []string                           `json:"artistsNames"`
}

type YoutubeSettings struct {
	MaxRequests             *int                      `json:"maxRequests"`
	AcceptOnlyWhenOnline    *bool                     `json:"acceptOnlyWhenOnline"`
	ChannelPointsRewardName *string                   `json:"channelPointsRewardName"`
	User                    *YoutubeUserSettings      `json:"user"`
	Song                    *YotubeSongSettings       `json:"song"`
	BlackList               *YoutubeBlacklistSettings `json:"blacklist"`
}
