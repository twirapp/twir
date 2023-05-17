package types

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"
	model "github.com/satont/tsuwari/libs/gomodels"
	"sync"
)

type ParseContextSender struct {
	ID          string
	Name        string
	DisplayName string
	Badges      []string
}

type ParseContextChannel struct {
	ID   string
	Name string
}

type ParseContextCacheLocks struct {
	Stream              *sync.Mutex
	DBUser              *sync.Mutex
	TwitchUser          *sync.Mutex
	TwitchFollow        *sync.Mutex
	TwitchChannel       *sync.Mutex
	ChannelIntegrations *sync.Mutex
	FaceitMatches       *sync.Mutex
	ValorantProfile     *sync.Mutex
	ValorantMatch       *sync.Mutex
}

type ParseContextCache struct {
	Stream          *model.ChannelsStreams
	DbUserStats     *model.UsersStats
	TwitchUser      *helix.User
	TwitchFollow    *helix.UserFollow
	TwitchChannel   *helix.ChannelInformation
	Integrations    []model.ChannelsIntegrations
	FaceitData      *FaceitResult
	ValorantProfile *ValorantProfile
	ValorantMatches []ValorantMatch

	Locks *ParseContextCacheLocks
}

type ParseContext struct {
	Channel *ParseContextChannel
	Sender  *ParseContextSender

	Text      *string
	IsCommand bool

	Services *services.Services
}
