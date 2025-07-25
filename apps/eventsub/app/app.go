package app

import (
	bus_listener "github.com/twirapp/twir/apps/eventsub/internal/bus-listener"
	"github.com/twirapp/twir/apps/eventsub/internal/handler"
	"github.com/twirapp/twir/apps/eventsub/internal/manager"
	"github.com/twirapp/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/baseapp"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	channelsongrequestssettingscache "github.com/twirapp/twir/libs/cache/channel_song_requests_settings"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	channelsintegrationssettingsseventvcache "github.com/twirapp/twir/libs/cache/channels_integrations_settings_seventv"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/websockets"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistpostgres "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistorypostgres "github.com/twirapp/twir/libs/repositories/channels_info_history/datasource/postgres"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	channelsredemptionshistoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_redemptions_history/datasources/clickhouse"

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
			channelseventslistpostgres.NewFx,
			fx.As(new(channelseventslist.Repository)),
		),
		fx.Annotate(
			alertsrepositorypgx.NewFx,
			fx.As(new(alertsrepository.Repository)),
		),
		fx.Annotate(
			channelsredemptionshistoryclickhouse.NewFx,
			fx.As(new(channelsredemptionshistory.Repository)),
		),
		channelcache.New,
		channelscommandsprefixcache.New,
		channelalertscache.New,
		commandscache.New,
		channelsongrequestssettingscache.New,
		channelsintegrationssettingsseventvcache.New,
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
