package cacher

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/types/services"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
)

type locks struct {
	stream      sync.Mutex
	dbUserStats sync.Mutex

	twitchSenderUser sync.Mutex
	twitchFollow     sync.Mutex
	twitchChannel    sync.Mutex

	channelIntegrations sync.Mutex

	faceitMatches  sync.Mutex
	faceitUserData sync.Mutex

	valorantProfile sync.Mutex
	valorantMatch   sync.Mutex
}

type cache struct {
	stream      *model.ChannelsStreams
	dbUserStats *model.UsersStats

	twitchSenderUser  *helix.User
	twitchUserFollows map[string]*helix.UserFollow
	twitchChannel     *helix.ChannelInformation

	channelIntegrations []*model.ChannelsIntegrations

	faceitData *types.FaceitResult

	valorantProfile *types.ValorantProfile
	valorantMatches []*types.ValorantMatch
}

type cacher struct {
	services        *services.Services
	parseCtxChannel *types.ParseContextChannel
	parseCtxSender  *types.ParseContextSender
	parseCtxText    *string
	cache           *cache
	locks           *locks
}

type CacherOpts struct {
	Services        *services.Services
	ParseCtxChannel *types.ParseContextChannel
	ParseCtxSender  *types.ParseContextSender
	ParseCtxText    *string
}

func NewCacher(opts *CacherOpts) types.DataCacher {
	return &cacher{
		services:        opts.Services,
		parseCtxChannel: opts.ParseCtxChannel,
		parseCtxSender:  opts.ParseCtxSender,
		parseCtxText:    opts.ParseCtxText,
		locks:           &locks{},
		cache: &cache{
			twitchUserFollows: make(map[string]*helix.UserFollow),
		},
	}
}

// GetChannelStream implements types.VariablesCacher
func (c *cacher) GetChannelStream(ctx context.Context) *model.ChannelsStreams {
	c.locks.stream.Lock()
	defer c.locks.stream.Unlock()

	if c.cache.stream != nil {
		return c.cache.stream
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil
	}

	stream := model.ChannelsStreams{}

	err = c.services.Gorm.
		WithContext(ctx).
		Where(`"userId" = ?`, c.parseCtxChannel.ID).
		First(&stream).Error

	if err != nil {
		c.services.Logger.Sugar().Error(err)

		streams, err := twitchClient.GetStreams(&helix.StreamsParams{
			UserIDs: []string{c.parseCtxChannel.ID},
		})

		if err != nil || len(streams.Data.Streams) == 0 {
			return nil
		}

		helixStream := streams.Data.Streams[0]

		tags := pq.StringArray{}
		for _, t := range helixStream.TagIDs {
			tags = append(tags, t)
		}
		stream = model.ChannelsStreams{
			ID:             helixStream.ID,
			UserId:         helixStream.UserID,
			UserLogin:      helixStream.UserLogin,
			UserName:       helixStream.UserName,
			GameId:         helixStream.GameID,
			GameName:       helixStream.GameName,
			CommunityIds:   []string{},
			Type:           helixStream.Type,
			Title:          helixStream.Title,
			ViewerCount:    helixStream.ViewerCount,
			StartedAt:      helixStream.StartedAt,
			Language:       helixStream.Language,
			ThumbnailUrl:   helixStream.ThumbnailURL,
			TagIds:         &tags,
			IsMature:       helixStream.IsMature,
			ParsedMessages: 0,
		}

		c.services.Gorm.Save(&stream)
		c.cache.stream = &stream
	} else {
		c.cache.stream = &stream
	}

	return c.cache.stream
}

// GetEnabledIntegrations implements types.VariablesCacher
func (c *cacher) GetEnabledChannelIntegrations(ctx context.Context) []*model.ChannelsIntegrations {
	c.locks.channelIntegrations.Lock()
	defer c.locks.channelIntegrations.Unlock()

	if c.cache.channelIntegrations != nil {
		return c.cache.channelIntegrations
	}

	var result []*model.ChannelsIntegrations
	err := c.services.Gorm.Where(`"channelId" = ? AND enabled = ?`, c.parseCtxChannel.ID, true).
		WithContext(ctx).
		Preload("Integration").
		Find(&result).
		Error

	if err == nil {
		c.cache.channelIntegrations = result
	}

	return c.cache.channelIntegrations
}

type faceitMatchesResponse []*types.FaceitMatch

// GetFaceitLatestMatches implements types.VariablesCacher
func (c *cacher) GetFaceitLatestMatches(ctx context.Context) ([]*types.FaceitMatch, error) {
	c.locks.faceitMatches.Lock()
	defer c.locks.faceitMatches.Unlock()

	if c.cache.faceitData == nil || c.cache.faceitData.FaceitUser == nil {
		faceitUser, err := c.GetFaceitUserData(ctx)
		if err != nil {
			return nil, err
		}

		c.cache.faceitData.FaceitUser = faceitUser
	}

	if c.cache.faceitData.Matches != nil {
		return c.cache.faceitData.Matches, nil
	}

	reqResult := faceitMatchesResponse{}

	_, err := req.C().EnableForceHTTP1().R().
		SetContext(ctx).
		SetSuccessResult(&reqResult).
		Get(fmt.Sprintf(
			"https://api.faceit.com/stats/api/v1/stats/time/users/%s/games/%s?size=30",
			c.cache.faceitData.FaceitUser.PlayerId,
			c.cache.faceitData.FaceitUser.FaceitGame.Name,
		))

	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil, err
	}

	var matches []*types.FaceitMatch
	stream := c.GetChannelStream(ctx)
	if stream == nil {
		return matches, nil
	}
	startedDate := stream.StartedAt.UnixMilli()

	for i, match := range reqResult {
		matchCreatedAt := time.UnixMilli(match.UpdateAt).UnixMilli()

		if matchCreatedAt < startedDate {
			continue
		}

		if i+1 > len(reqResult)-1 {
			break
		}

		val := false
		if match.RawIsWin == "1" {
			val = true
		}
		match.IsWin = val

		if i+1 >= len(reqResult)-1 {
			break
		}

		prevMatch := reqResult[i+1]
		if prevMatch == nil || prevMatch.Elo == nil || match.Elo == nil {
			continue
		}

		prevElo, pErr := strconv.Atoi(*prevMatch.Elo)
		currElo, cErr := strconv.Atoi(*match.Elo)

		if pErr != nil || cErr != nil {
			continue
		}

		var eloDiff int
		if *prevMatch.Elo > *match.Elo {
			eloDiff = -(prevElo - currElo)
		} else {
			eloDiff = currElo - prevElo
		}

		newMatchEloDiff := strconv.Itoa(eloDiff)
		match.EloDiff = &newMatchEloDiff
		matches = append(matches, match)
	}

	c.cache.faceitData.Matches = matches

	return matches, nil
}

// GetFaceitTodayEloDiff implements types.VariablesCacher
func (c *cacher) GetFaceitTodayEloDiff(ctx context.Context, matches []*types.FaceitMatch) int {
	if matches == nil {
		return 0
	}

	sum := lo.Reduce(matches, func(agg int, item *types.FaceitMatch, _ int) int {
		if item.EloDiff == nil {
			return agg
		}
		v, err := strconv.Atoi(*item.EloDiff)
		if err != nil {
			return agg
		}
		return agg + v
	}, 0)

	return sum
}

// GetFaceitUserData implements types.VariablesCacher
func (c *cacher) GetFaceitUserData(ctx context.Context) (*types.FaceitUser, error) {
	c.locks.faceitUserData.Lock()
	defer c.locks.faceitUserData.Unlock()

	if c.cache.faceitData != nil && c.cache.faceitData.FaceitUser != nil {
		return c.cache.faceitData.FaceitUser, nil
	}

	c.cache.faceitData = &types.FaceitResult{}

	integrations := c.GetEnabledChannelIntegrations(ctx)

	if integrations == nil {
		return nil, errors.New("no enabled integrations")
	}

	integration, ok := lo.Find(integrations, func(i *model.ChannelsIntegrations) bool {
		return i.Integration.Service == "FACEIT" && i.Enabled
	})

	if !ok {
		return nil, errors.New("faceit integration not enabled")
	}

	var game string = *integration.Data.Game

	if integration.Data.Game == nil {
		game = "csgo"
	}

	data := &types.FaceitUserResponse{}
	resp, err := req.C().EnableForceHTTP1().R().
		SetContext(ctx).
		SetBearerAuthToken(integration.Integration.APIKey.String).
		SetSuccessResult(data).
		Get("https://open.faceit.com/data/v4/players/" + *integration.Data.UserId)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New(
			"user not found on faceit. Please make sure you typed correct nickname",
		)
	}

	if data.Games[game] == nil {
		return nil, errors.New(game + " game not found in faceit response.")
	}

	data.Games[game].Name = game

	c.cache.faceitData.FaceitUser = &types.FaceitUser{
		FaceitGame: *data.Games[game],
		PlayerId:   data.PlayerId,
	}

	return c.cache.faceitData.FaceitUser, nil
}

// GetFollowAge implements types.VariablesCacher
func (c *cacher) GetTwitchUserFollow(ctx context.Context, userID string) *helix.UserFollow {
	c.locks.twitchFollow.Lock()
	defer c.locks.twitchFollow.Unlock()

	if c.cache.twitchUserFollows[userID] != nil {
		return c.cache.twitchUserFollows[userID]
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil
	}

	follow, err := twitchClient.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: userID,
		ToID:   c.parseCtxChannel.ID,
	})

	if err == nil && len(follow.Data.Follows) != 0 {
		c.cache.twitchUserFollows[userID] = &follow.Data.Follows[0]
	}

	return c.cache.twitchUserFollows[userID]
}

// GetGbUser implements types.VariablesCacher
func (c *cacher) GetGbUserStats(ctx context.Context) *model.UsersStats {
	c.locks.dbUserStats.Lock()
	defer c.locks.dbUserStats.Unlock()

	if c.cache.dbUserStats != nil {
		return c.cache.dbUserStats
	}

	result := &model.UsersStats{}

	err := c.services.Gorm.
		WithContext(ctx).
		Where(`"userId" = ? AND "channelId" = ?`, c.parseCtxSender.ID, c.parseCtxChannel.ID).
		Find(result).
		Error
	if err == nil {
		c.cache.dbUserStats = result
	}

	return c.cache.dbUserStats
}

// GetTwitchChannel implements types.VariablesCacher
func (c *cacher) GetTwitchChannel(ctx context.Context) *helix.ChannelInformation {
	c.locks.twitchChannel.Lock()
	defer c.locks.twitchChannel.Unlock()

	if c.cache.twitchChannel != nil {
		return c.cache.twitchChannel
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	channel, err := twitchClient.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{c.parseCtxChannel.ID},
	})

	if err == nil && len(channel.Data.Channels) != 0 {
		c.cache.twitchChannel = &channel.Data.Channels[0]
	}

	return c.cache.twitchChannel
}

// GetTwitchUser implements types.VariablesCacher
func (c *cacher) GetTwitchSenderUser(ctx context.Context) *helix.User {
	c.locks.twitchSenderUser.Lock()
	defer c.locks.twitchSenderUser.Unlock()

	if c.cache.twitchSenderUser != nil {
		return c.cache.twitchSenderUser
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil
	}

	users, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: []string{c.parseCtxSender.ID},
	})

	if err == nil && len(users.Data.Users) != 0 {
		c.cache.twitchSenderUser = &users.Data.Users[0]
	}

	return c.cache.twitchSenderUser
}

// GetValorantMatches implements types.VariablesCacher
func (c *cacher) GetValorantMatches(ctx context.Context) []*types.ValorantMatch {
	c.locks.valorantMatch.Lock()
	defer c.locks.valorantMatch.Unlock()

	if c.cache.valorantMatches != nil {
		return c.cache.valorantMatches
	}

	var data *types.ValorantMatchesResponse

	integrations := c.GetEnabledChannelIntegrations(ctx)
	integration, ok := lo.Find(integrations, func(item *model.ChannelsIntegrations) bool {
		return item.Integration.Service == "VALORANT"
	})

	if !ok || integration.Data == nil || integration.Data.UserName == nil {
		return nil
	}

	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get("https://api.henrikdev.xyz/valorant/v3/matches/eu/" + strings.Replace(
			*integration.Data.UserName,
			"#",
			"/",
			1,
		))
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	c.cache.valorantMatches = data.Data

	return c.cache.valorantMatches
}

// GetValorantProfile implements types.VariablesCacher
func (c *cacher) GetValorantProfile(ctx context.Context) *types.ValorantProfile {
	c.locks.valorantProfile.Lock()
	defer c.locks.valorantProfile.Unlock()

	if c.cache.valorantProfile != nil {
		return c.cache.valorantProfile
	}

	integrations := c.GetEnabledChannelIntegrations(ctx)
	integration, ok := lo.Find(integrations, func(item *model.ChannelsIntegrations) bool {
		return item.Integration.Service == "VALORANT"
	})

	if !ok || integration.Data == nil || integration.Data.UserName == nil {
		return nil
	}

	c.cache.valorantProfile = &types.ValorantProfile{}
	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(c.cache.valorantProfile).
		Get("https://api.henrikdev.xyz/valorant/v1/mmr/eu/" + strings.Replace(
			*integration.Data.UserName,
			"#",
			"/",
			1,
		))
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	return c.cache.valorantProfile
}
