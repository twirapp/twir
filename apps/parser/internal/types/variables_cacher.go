package types

import (
	"context"

	"github.com/nicklaw5/helix/v2"
	model "github.com/twirapp/twir/libs/gomodels"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	"github.com/twirapp/twir/libs/twitch"
)

type DataCacher interface {
	GetEnabledChannelIntegrations(ctx context.Context) []*model.ChannelsIntegrations

	GetFaceitLatestMatches(ctx context.Context) ([]*FaceitMatch, error)
	GetFaceitTodayEloDiff(ctx context.Context, matches []*FaceitMatch) int
	GetFaceitUserData(ctx context.Context) (*FaceitUser, error)
	ComputeFaceitGainLoseEstimate(ctx context.Context) (*FaceitEstimateGainLose, error)

	GetTwitchUserFollow(ctx context.Context, userId string) *helix.ChannelFollow
	GetGbUserStats(ctx context.Context, userId string) *model.UsersStats
	GetTwitchChannel(ctx context.Context) *helix.ChannelInformation
	GetTwitchSenderUser(ctx context.Context) *helix.User
	GetTwitchUserById(ctx context.Context, userId string) (*helix.User, error)
	GetTwitchUserByName(ctx context.Context, userId string) (*helix.User, error)
	GetValorantMatches(ctx context.Context) []ValorantMatch
	GetValorantMMR(ctx context.Context) *ValorantMMR

	GetCurrentSong(ctx context.Context) *CurrentSong

	GetSeventvProfileGetTwitchId(
		ctx context.Context,
		userId string,
	) (*seventvintegrationapi.TwirSeventvUser, error)

	GetSubAgeInfo(ctx context.Context, channelName, userName string) (
		*twitch.UserSubscribePayload,
		error,
	)
}
