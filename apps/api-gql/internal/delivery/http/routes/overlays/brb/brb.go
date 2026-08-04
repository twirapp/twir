package brb

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	buscore "github.com/twirapp/twir/libs/bus-core"
)

type Dependencies struct {
	API            huma.API
	Service        *be_right_back.Service
	TwirBus        *buscore.Bus
	ChannelService *channels.Service
	Middlewares    *middlewares.Middlewares
	Sessions       *auth.Auth
}

type Registration struct{}

func RegisterRoutes(deps Dependencies) Registration {
	opts := func() StartOpts {
		return StartOpts{
			Service: deps.Service, TwirBus: deps.TwirBus, ChannelService: deps.ChannelService,
			Middlewares: deps.Middlewares, Sessions: deps.Sessions,
		}
	}
	newStart(opts()).Register(deps.API)
	startOpts := opts()
	newStop(StopOpts(startOpts)).Register(deps.API)
	return Registration{}
}
