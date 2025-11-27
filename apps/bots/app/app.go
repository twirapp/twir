package app

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	bus_listener "github.com/twirapp/twir/apps/bots/internal/bus-listener"
	"github.com/twirapp/twir/apps/bots/internal/messagehandler"
	mod_task_queue "github.com/twirapp/twir/apps/bots/internal/mod-task-queue"
	"github.com/twirapp/twir/apps/bots/internal/moderationhelpers"
	"github.com/twirapp/twir/apps/bots/internal/services/keywords"
	toxicity_check "github.com/twirapp/twir/apps/bots/internal/services/toxicity-check"
	"github.com/twirapp/twir/apps/bots/internal/services/tts"
	stream_handlers "github.com/twirapp/twir/apps/bots/internal/stream-handlers"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	"github.com/twirapp/twir/apps/bots/pkg/tlds"
	"github.com/twirapp/twir/libs/baseapp"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	channelsmoderationsettingscache "github.com/twirapp/twir/libs/cache/channels_moderation_settings"
	chatwallcacher "github.com/twirapp/twir/libs/cache/chat_wall"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	giveawayscache "github.com/twirapp/twir/libs/cache/giveaways"
	greetingscache "github.com/twirapp/twir/libs/cache/greetings"
	keywordscache "github.com/twirapp/twir/libs/cache/keywords"
	rolescache "github.com/twirapp/twir/libs/cache/roles"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsemotesusagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/datasources/clickhouse"
	channelsmoderationsettingsrepository "github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	channelsmoderationsettingsrepositorypostgres "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/datasource/postgres"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallrepositorypostgres "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"
	giveawaysrepository "github.com/twirapp/twir/libs/repositories/giveaways"
	giveawaysrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways/pgx"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
	keywordsrepositorypgx "github.com/twirapp/twir/libs/repositories/keywords/pgx"
	overlays_tts_repository "github.com/twirapp/twir/libs/repositories/overlays_tts"
	overlays_tts_pgx "github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
	rolesrepository "github.com/twirapp/twir/libs/repositories/roles"
	rolesrepositorypgx "github.com/twirapp/twir/libs/repositories/roles/pgx"
	sentmessagesrepository "github.com/twirapp/twir/libs/repositories/sentmessages"
	sentmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/sentmessages/pgx"
	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	toxicmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/toxic_messages/pgx"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	usersstatsrepository "github.com/twirapp/twir/libs/repositories/users_stats"
	usersstatsrepositorypostgres "github.com/twirapp/twir/libs/repositories/users_stats/datasources/postgres"
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
		fx.Annotate(
			chatmessagesrepositoryclickhouse.NewFx,
			fx.As(new(chatmessagesrepository.Repository)),
		),
		fx.Annotate(
			channelscommandsprefixpgx.NewFx,
			fx.As(new(channelscommandsprefixrepository.Repository)),
		),
		fx.Annotate(
			chatwallrepositorypostgres.NewFx,
			fx.As(new(chatwallrepository.Repository)),
		),
		fx.Annotate(
			giveawaysrepositorypgx.NewFx,
			fx.As(new(giveawaysrepository.Repository)),
		),
		fx.Annotate(
			channelsmoderationsettingsrepositorypostgres.NewFx,
			fx.As(new(channelsmoderationsettingsrepository.Repository)),
		),
		fx.Annotate(
			channelsemotesusagesrepositoryclickhouse.NewFx,
			fx.As(new(channelsemotesusagesrepository.Repository)),
		),
		fx.Annotate(
			usersrepositorypgx.NewFx,
			fx.As(new(usersrepository.Repository)),
		),
		fx.Annotate(
			usersstatsrepositorypostgres.NewFx,
			fx.As(new(usersstatsrepository.Repository)),
		),
		fx.Annotate(
			rolesrepositorypgx.NewFx,
			fx.As(new(rolesrepository.Repository)),
		),
		fx.Annotate(
			overlays_tts_pgx.NewFx,
			fx.As(new(overlays_tts_repository.Repository)),
		),
	),
	fx.Provide(
		tlds.New,
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		workers.New,
		chatwallcacher.NewEnabledOnly,
		chatwallcacher.NewSettings,
		giveawayscache.New,
		fx.Annotate(
			mod_task_queue.NewRedisModTaskDistributor,
			fx.As(new(mod_task_queue.TaskDistributor)),
		),
		rolescache.New,
		toxicity_check.New,
		func(
			repo channelscommandsprefixrepository.Repository,
			bus *buscore.Bus,
		) *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix] {
			return channelscommandsprefixcache.New(repo, bus)
		},
		ttscache.NewTTSSettings,
		keywordscache.New,
		greetingscache.New,
		channelcache.New,
		twitchactions.New,
		channelsmoderationsettingscache.New,
		moderationhelpers.New,
		messagehandler.New,
		keywords.New,
		tts.New,
	),
	fx.Invoke(
		mod_task_queue.NewRedisTaskProcessor,
		func(config cfg.Config) {
			if config.AppEnv != "development" {
				http.Handle("/metrics", promhttp.Handler())
				go http.ListenAndServe("0.0.0.0:3000", nil)
			}
		},
		stream_handlers.New,
		bus_listener.New,
		func(l *slog.Logger) {
			l.Info("🚀 Bots started")
		},
	),
)
