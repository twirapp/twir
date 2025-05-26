package app

import (
	eventsActivity "github.com/satont/twir/apps/events/internal/activities/events"
	"github.com/satont/twir/apps/events/internal/chat_alerts"
	"github.com/satont/twir/apps/events/internal/hydrator"
	"github.com/satont/twir/apps/events/internal/listener"
	"github.com/satont/twir/apps/events/internal/song_request"
	"github.com/satont/twir/apps/events/internal/workers"
	"github.com/satont/twir/apps/events/internal/workflows"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	channelseventswithoperations "github.com/twirapp/twir/libs/cache/channels_events_with_operations"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/grpc/ytsr"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"events",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "events"}),
	fx.Provide(
		fx.Annotate(
			greetingsrepositorypgx.NewFx,
			fx.As(new(greetingsrepository.Repository)),
		),
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		func(config cfg.Config) ytsr.YtsrClient {
			return clients.NewYtsr(config.AppEnv)
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
			l.Info("Events service started")
		},
	),
)
