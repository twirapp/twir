package modules

import "github.com/satont/tsuwari/libs/types/types"

type YouTube struct {
	POST                    types.YoutubeBlacklistSettingsChannels
	GET                     types.YoutubeBlacklistSettings
	POST_BLACKLIST_SONGS    types.YoutubeBlacklistSettingsSongs
	POST_BLACKLIST_CHANNELS types.YoutubeBlacklistSettingsChannels
	POST_BLACKLIST_USERS    types.YoutubeBlacklistSettingsUsers
}
