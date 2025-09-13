package app

import (
	eventsActivity "github.com/twirapp/twir/apps/events/internal/activities/events"
	"github.com/twirapp/twir/apps/events/internal/chat_alerts"
	"github.com/twirapp/twir/apps/events/internal/hydrator"
	"github.com/twirapp/twir/apps/events/internal/listener"
	"github.com/twirapp/twir/apps/events/internal/song_request"
	"github.com/twirapp/twir/apps/events/internal/workers"
	"github.com/twirapp/twir/apps/events/internal/workflows"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/cache/channel"
	channelseventswithoperations "github.com/twirapp/twir/libs/cache/channels_events_with_operations"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistpostgres "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts"
	channelsmodules_settingsttspgx "github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts/postgres"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"

	eventsrepository "github.com/twirapp/twir/libs/repositories/events"
	eventsrepositorypostgres "github.com/twirapp/twir/libs/repositories/events/pgx"

	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypostgres "github.com/twirapp/twir/libs/repositories/channels/pgx"

	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	variablesrepositorypostgres "github.com/twirapp/twir/libs/repositories/variables/pgx"
)

var App = fx.Module(
	"events",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "events"}),
	fx.Provide(
		fx.Annotate(
			greetingsrepositorypgx.NewFx,
			fx.As(new(greetingsrepository.Repository)),
		),
		fx.Annotate(
			channelsrepositorypostgres.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
		fx.Annotate(
			eventsrepositorypostgres.NewFx,
			fx.As(new(eventsrepository.Repository)),
		),
		fx.Annotate(
			variablesrepositorypostgres.NewFx,
			fx.As(new(variablesrepository.Repository)),
		),
		fx.Annotate(
			channelseventslistpostgres.NewFx,
			fx.As(new(channelseventslist.Repository)),
		),
		fx.Annotate(
			channelsmodules_settingsttspgx.NewFx,
			fx.As(new(channels_modules_settings_tts.Repository)),
		),
		channel.New,
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		song_request.New,
		hydrator.New,
		eventsActivity.New,
		workflows.NewEventsWorkflow,
		chat_alerts.New,
		channelseventswithoperations.New,
		chatalertscache.New,
	),
	fx.Invoke(
		uptrace.NewFx("events"),
		workers.NewEventsWorker,
		listener.New,
		func(l logger.Logger) {
			l.Info("🤖 Events service started")
		},
	),
)
