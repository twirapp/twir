package app

import (
	bus_listener "github.com/satont/twir/apps/eventsub/internal/bus-listener"
	"github.com/satont/twir/apps/eventsub/internal/handler"
	"github.com/satont/twir/apps/eventsub/internal/manager"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/baseapp"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	channelsongrequestssettingscache "github.com/twirapp/twir/libs/cache/channel_song_requests_settings"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	channelredemptionshistory "github.com/twirapp/twir/libs/repositories/channel_redemptions_history"
	channelredemptionshistorypostgres "github.com/twirapp/twir/libs/repositories/channel_redemptions_history/datasources/postgres"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistpostgres "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistorypostgres "github.com/twirapp/twir/libs/repositories/channels_info_history/datasource/postgres"

	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"

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
		fx.Annotate(
			channelsinfohistorypostgres.NewFx,
			fx.As(new(channelsinfohistory.Repository)),
		),
		fx.Annotate(
			streamsrepositorypostgres.NewFx,
			fx.As(new(streamsrepository.Repository)),
		),
		fx.Annotate(
			channelredemptionshistorypostgres.NewFx,
			fx.As(new(channelredemptionshistory.Repository)),
		),
		fx.Annotate(
			channelseventslistpostgres.NewFx,
			fx.As(new(channelseventslist.Repository)),
		),
		fx.Annotate(
			alertsrepositorypgx.NewFx,
			fx.As(new(alertsrepository.Repository)),
		),
		channelcache.New,
		channelscommandsprefixcache.New,
		channelalertscache.New,
		commandscache.New,
		channelsongrequestssettingscache.New,
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
