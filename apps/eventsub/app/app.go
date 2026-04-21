package app

import (
	"context"
	"log/slog"

	bus_listener "github.com/twirapp/twir/apps/eventsub/internal/bus-listener"
	"github.com/twirapp/twir/apps/eventsub/internal/handler"
	httpserver "github.com/twirapp/twir/apps/eventsub/internal/http"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	"github.com/twirapp/twir/apps/eventsub/internal/manager"
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	"github.com/twirapp/twir/apps/eventsub/internal/webhook"
	"github.com/twirapp/twir/libs/baseapp"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	channelsongrequestssettingscache "github.com/twirapp/twir/libs/cache/channel_song_requests_settings"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	channelsintegrationssettingsseventvcache "github.com/twirapp/twir/libs/cache/channels_integrations_settings_seventv"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/otel"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistpostgres "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistorypostgres "github.com/twirapp/twir/libs/repositories/channels_info_history/datasource/postgres"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	channelsredemptionshistoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_redemptions_history/datasources/clickhouse"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsespostgres "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	usersstats "github.com/twirapp/twir/libs/repositories/users_stats"
	usersstatsrepositorypostgres "github.com/twirapp/twir/libs/repositories/users_stats/datasources/postgres"
	userswithstatsrepository "github.com/twirapp/twir/libs/repositories/userswithstats"
	userswithstatsrepositorypostgres "github.com/twirapp/twir/libs/repositories/userswithstats/datasource/postgres"
	"go.uber.org/fx"

	twitchconduitsrepository "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	twitchconduitsrepositorypostgres "github.com/twirapp/twir/libs/repositories/twitch_conduits/datasource/postgres"
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
		fx.Annotate(
			twitchconduitsrepositorypostgres.NewFx,
			fx.As(new(twitchconduitsrepository.Repository)),
		),
		fx.Annotate(
			commandswithgroupsandresponsespostgres.NewFx,
			fx.As(new(commandswithgroupsandresponsesrepository.Repository)),
		),
		fx.Annotate(
			userswithstatsrepositorypostgres.NewFx,
			fx.As(new(userswithstatsrepository.Repository)),
		),
		fx.Annotate(
			usersstatsrepositorypostgres.NewFx,
			fx.As(new(usersstats.Repository)),
		),
		fx.Annotate(
			usersrepositorypgx.NewFx,
			fx.As(new(usersrepository.Repository)),
		),
		fx.Annotate(
			kickbotsrepositorypgx.NewFx,
			fx.As(new(kickbotsrepository.Repository)),
		),
		channelcache.New,
		func(
			repo channelscommandsprefixrepository.Repository,
			bus *buscore.Bus,
		) *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix] {
			return channelscommandsprefixcache.New(repo, bus)
		},
		user_creator.New,
		channelalertscache.New,
		commandscache.New,
		channelsongrequestssettingscache.New,
		channelsintegrationssettingsseventvcache.New,
		manager.NewManager,
		handler.New,
		httpserver.New,
		kick.New,
		kick.NewHandlers,
		kick.NewResubscribeJob,
		webhook.NewManager,
	),
	fx.Invoke(
		otel.NewFx("eventsub"),
		bus_listener.New,
		func(s *httpserver.Server, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error { return s.Start() },
				OnStop:  func(ctx context.Context) error { return s.Stop(ctx) },
			})
		},
		func(l *slog.Logger) {
			l.Info("🚀 EventSub App started")
		},
	),
)
