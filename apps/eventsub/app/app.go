package app

import (
	bus_listener "github.com/satont/twir/apps/eventsub/internal/bus-listener"
	"github.com/satont/twir/apps/eventsub/internal/handler"
	"github.com/satont/twir/apps/eventsub/internal/manager"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/baseapp"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"

	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Options(
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "eventsub"}),
	fx.Provide(
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		func(config cfg.Config) events.EventsClient {
			return clients.NewEvents(config.AppEnv)
		},
		func(config cfg.Config) parser.ParserClient {
			return clients.NewParser(config.AppEnv)
		},
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		fx.Annotate(
			channelsrepositorypgx.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
		fx.Annotate(
			channelscommandsprefixpgx.NewFx,
			fx.As(new(channelscommandsprefixrepository.Repository)),
		),
		fx.Annotate(
			scheduledvipsrepositorypgx.NewFx,
			fx.As(new(scheduledvipsrepository.Repository)),
		),
		channelcache.New,
		channelscommandsprefixcache.New,
		tunnel.New,
		manager.NewCreds,
		manager.NewManager,
		handler.New,
	),
	fx.Invoke(
		uptrace.NewFx("eventsub"),
		handler.New,
		bus_listener.New,
	),
)
