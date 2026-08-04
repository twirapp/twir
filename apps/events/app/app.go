package app

import (
	"fmt"
	"log/slog"

	"github.com/goforj/wire"
	eventsactivity "github.com/twirapp/twir/apps/events/internal/activities/events"
	"github.com/twirapp/twir/apps/events/internal/chat_alerts"
	"github.com/twirapp/twir/apps/events/internal/hydrator"
	"github.com/twirapp/twir/apps/events/internal/listener"
	"github.com/twirapp/twir/apps/events/internal/song_request"
	"github.com/twirapp/twir/apps/events/internal/workers"
	"github.com/twirapp/twir/apps/events/internal/workflows"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelseventswithoperations "github.com/twirapp/twir/libs/cache/channels_events_with_operations"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	"github.com/twirapp/twir/libs/cache/tts"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	grpcwebsockets "github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelseventsrepository "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	commandsrepository "github.com/twirapp/twir/libs/repositories/commands"
	commandsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands/pgx"
	eventsrepository "github.com/twirapp/twir/libs/repositories/events"
	eventsrepositorypgx "github.com/twirapp/twir/libs/repositories/events/pgx"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	overlaysttsrepository "github.com/twirapp/twir/libs/repositories/overlays_tts"
	overlaysttsrepositorypgx "github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypgx "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	variablesrepositorypgx "github.com/twirapp/twir/libs/repositories/variables/pgx"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

const Service = "events"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	greetingsrepositorypgx.NewFx,
	wire.Bind(new(greetingsrepository.Repository), new(*greetingsrepositorypgx.Pgx)),
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	streamsrepositorypgx.NewFx,
	wire.Bind(new(streamsrepository.Repository), new(*streamsrepositorypgx.Pgx)),
	commandsrepositorypgx.NewFx,
	wire.Bind(new(commandsrepository.Repository), new(*commandsrepositorypgx.Pgx)),
	eventsrepositorypgx.NewFx,
	wire.Bind(new(eventsrepository.Repository), new(*eventsrepositorypgx.Pgx)),
	variablesrepositorypgx.NewFx,
	wire.Bind(new(variablesrepository.Repository), new(*variablesrepositorypgx.Pgx)),
	channelseventsrepositorypgx.NewFx,
	wire.Bind(new(channelseventsrepository.Repository), new(*channelseventsrepositorypgx.Pgx)),
	overlaysttsrepositorypgx.NewFx,
	wire.Bind(new(overlaysttsrepository.Repository), new(*overlaysttsrepositorypgx.Pgx)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	channelcache.New,
	channelservice.NewChannelService,
	tts.NewTTSSettings,
	NewWebsocketClient,
	wire.Struct(new(song_request.Opts), "*"),
	song_request.New,
	wire.Struct(new(hydrator.Opts), "*"),
	hydrator.New,
	wire.Struct(new(eventsactivity.Opts), "*"),
	eventsactivity.New,
	wire.Struct(new(workflows.EventsWorkflowOpts), "*"),
	workflows.NewEventsWorkflow,
	wire.Struct(new(chat_alerts.Opts), "*"),
	chat_alerts.New,
	channelseventswithoperations.New,
	chatalertscache.New,
	wire.Struct(new(workers.EventsWorkerOpts), "*"),
	wire.Struct(new(listener.Opts), "*"),
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewWebsocketClient(config config.Config) grpcwebsockets.WebsocketClient {
	return clients.NewWebsocket(config.AppEnv)
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	workerOpts workers.EventsWorkerOpts,
	listenerOpts listener.Opts,
) (*Application, error) {
	if err := workers.NewEventsWorker(workerOpts); err != nil {
		return nil, fmt.Errorf("create events worker: %w", err)
	}
	if err := listener.New(listenerOpts); err != nil {
		return nil, fmt.Errorf("create events listener: %w", err)
	}

	logger.Info("🤖 Events service started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
