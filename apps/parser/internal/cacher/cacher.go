package cacher

import (
	"context"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
)

type locks struct {
	stream      sync.Mutex
	dbUserStats sync.Mutex

	twitchFollow            sync.Mutex
	twitchChannel           sync.Mutex
	cachedTwitchUsersById   sync.Mutex
	cachedTwitchUsersByName sync.Mutex
	channelIntegrations     sync.Mutex
	faceitMatches           sync.Mutex
	faceitUserData          sync.Mutex
	valorantProfile         sync.Mutex
	valorantMatches         sync.Mutex
	currentSong             sync.Mutex
	subage                  sync.Mutex
}

type cache struct {
	stream      *model.ChannelsStreams
	dbUserStats *model.UsersStats

	twitchUserFollows       map[string]*helix.ChannelFollow
	twitchChannel           *helix.ChannelInformation
	cachedTwitchUsersById   map[string]*helix.User
	cachedTwitchUsersByName map[string]*helix.User

	faceitData *types.FaceitResult

	valorantProfile *types.ValorantProfile

	currentSong *types.CurrentSong

	channelIntegrations []*model.ChannelsIntegrations

	valorantMatches []types.ValorantMatch

	seventvprofile *seventvintegrationapi.TwirSeventvUser

	cachedSubAgeInfo *twitch.UserSubscribePayload
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
			twitchUserFollows:       make(map[string]*helix.ChannelFollow),
			cachedTwitchUsersById:   make(map[string]*helix.User),
			cachedTwitchUsersByName: make(map[string]*helix.User),
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
