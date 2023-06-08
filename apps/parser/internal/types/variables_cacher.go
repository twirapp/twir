package types

import (
	"context"

	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

type DataCacher interface {
	GetChannelStream(ctx context.Context) *model.ChannelsStreams
	GetEnabledChannelIntegrations(ctx context.Context) []*model.ChannelsIntegrations
	GetFaceitLatestMatches(ctx context.Context) ([]*FaceitMatch, error)
	GetFaceitTodayEloDiff(ctx context.Context, matches []*FaceitMatch) int
	GetFaceitUserData(ctx context.Context) (*FaceitUser, error)
	GetRocketLeagueUserStats(ctx context.Context) *RocketLeagueUserStats
	GetTwitchUserFollow(ctx context.Context, userId string) *helix.UserFollow
	GetGbUserStats(ctx context.Context) *model.UsersStats
	GetTwitchChannel(ctx context.Context) *helix.ChannelInformation
	GetTwitchSenderUser(ctx context.Context) *helix.User
	GetValorantMatches(ctx context.Context) []*ValorantMatch
	GetValorantProfile(ctx context.Context) *ValorantProfile
}
