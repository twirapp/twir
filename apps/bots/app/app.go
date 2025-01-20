package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	bus_listener "github.com/satont/twir/apps/bots/internal/bus-listener"
	"github.com/satont/twir/apps/bots/internal/messagehandler"
	mod_task_queue "github.com/satont/twir/apps/bots/internal/mod-task-queue"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/services/keywords"
	toxicity_check "github.com/satont/twir/apps/bots/internal/services/toxicity-check"
	stream_handlers "github.com/satont/twir/apps/bots/internal/stream-handlers"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	keywordscache "github.com/twirapp/twir/libs/cache/keywords"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
	keywordsrepositorypgx "github.com/twirapp/twir/libs/repositories/keywords/pgx"
	sentmessagesrepository "github.com/twirapp/twir/libs/repositories/sentmessages"
	sentmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/sentmessages/pgx"
	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	toxicmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/toxic_messages/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"bots",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "bots"}),
	// repositories
	fx.Provide(
		fx.Annotate(
			keywordsrepositorypgx.NewFx,
			fx.As(new(keywordsrepository.Repository)),
		),
		fx.Annotate(
			greetingsrepositorypgx.NewFx,
			fx.As(new(greetingsrepository.Repository)),
		),
		fx.Annotate(
			sentmessagesrepositorypgx.NewFx,
			fx.As(new(sentmessagesrepository.Repository)),
		),
		fx.Annotate(
			channelsrepositorypgx.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
		fx.Annotate(
			toxicmessagesrepositorypgx.NewFx,
			fx.As(new(toxicmessagesrepository.Repository)),
		),
	),
	fx.Provide(
		tlds.New,
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
			mod_task_queue.NewRedisModTaskDistributor,
			fx.As(new(mod_task_queue.TaskDistributor)),
		),
		toxicity_check.New,
		keywordscache.New,
		twitchactions.New,
		moderationhelpers.New,
		messagehandler.New,
		keywords.New,
	),
	fx.Invoke(
		uptrace.NewFx("bots"),
		func(config cfg.Config) {
			if config.AppEnv != "development" {
				http.Handle("/metrics", promhttp.Handler())
				go http.ListenAndServe("0.0.0.0:3000", nil)
			}
		},
		stream_handlers.New,
		bus_listener.New,
		func(l logger.Logger) {
			l.Info("Bots started")
		},
		mod_task_queue.NewRedisTaskProcessor,
	),
)
