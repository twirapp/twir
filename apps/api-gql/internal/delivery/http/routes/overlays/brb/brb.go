package brb

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	buscore "github.com/twirapp/twir/libs/bus-core"
)

type Registration struct{}

func RegisterRoutes(api huma.API, service *be_right_back.Service, twirBus *buscore.Bus, channelService *channels.Service, middlewaresService *middlewares.Middlewares, sessions *auth.Auth) Registration {
	opts := func() StartOpts {
		return StartOpts{
			Service: service, TwirBus: twirBus, ChannelService: channelService,
			Middlewares: middlewaresService, Sessions: sessions,
		}
	}
	newStart(opts()).Register(api)
	startOpts := opts()
	newStop(StopOpts(startOpts)).Register(api)
	return Registration{}
}
