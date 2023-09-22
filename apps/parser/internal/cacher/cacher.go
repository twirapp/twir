package cacher

import (
	"context"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	model "github.com/satont/twir/libs/gomodels"
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

	currentSong sync.Mutex
}

type cache struct {
	stream      *model.ChannelsStreams
	dbUserStats *model.UsersStats

	twitchSenderUser  *helix.User
	twitchUserFollows map[string]*helix.ChannelFollow
	twitchChannel     *helix.ChannelInformation

	channelIntegrations []*model.ChannelsIntegrations

	faceitData *types.FaceitResult

	valorantProfile *types.ValorantProfile
	valorantMatches []*types.ValorantMatch

	currentSong *types.CurrentSong
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
			twitchUserFollows: make(map[string]*helix.ChannelFollow),
		},
	}
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
