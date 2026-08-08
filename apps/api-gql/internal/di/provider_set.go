package di

import (
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/goforj/wire"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/app"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	publicroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public"
	v2publicroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public/v2"
	http_webhooks "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-webhooks"
	httpmiddlewares "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	channelsfilesroute "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/channels/channels_files"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/integrations/valorant"
	mcpOAuthRoutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/mcp_oauth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/overlays/brb"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/pastebins"
	scheduledvipsroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/scheduled_vips"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/shortlinks"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/stream"
	ttsroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/tts"
	uploaderoutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/uploader"
	mcpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	twitchplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/platform/vkvideo"
	youtubeplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/youtube"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_commands_prefix"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_emotes_usages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_files"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_moderation_settings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_redemptions_history"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_secret"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_storage"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_messages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_translation"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/community_redemptions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/dashboard"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
	dashboard_widgets "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widgets"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	"github.com/twirapp/twir/apps/api-gql/internal/services/discord_integration"
	donatellointegration "github.com/twirapp/twir/apps/api-gql/internal/services/donatello_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/donatepay_integration"
	donatestreamintegration "github.com/twirapp/twir/apps/api-gql/internal/services/donatestream_integration"
	donationalertsintegration "github.com/twirapp/twir/apps/api-gql/internal/services/donationalerts_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/events"
	faceitintegration "github.com/twirapp/twir/apps/api-gql/internal/services/faceit_integration"
	gamesvoteban "github.com/twirapp/twir/apps/api-gql/internal/services/games_voteban"
	"github.com/twirapp/twir/apps/api-gql/internal/services/giveaways"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	lastfmintegration "github.com/twirapp/twir/apps/api-gql/internal/services/lastfm_integration"
	mcpOAuthService "github.com/twirapp/twir/apps/api-gql/internal/services/mcp_oauth"
	nightbotintegration "github.com/twirapp/twir/apps/api-gql/internal/services/nightbot_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/obs_websocket_module"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays_dudes"
	pastebinsservice "github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/services/quotes"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_with_roles_users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/scheduledvips"
	"github.com/twirapp/twir/apps/api-gql/internal/services/seventv_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	songrequestoverlaysettings "github.com/twirapp/twir/apps/api-gql/internal/services/song_request_overlay_settings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/song_requests"
	"github.com/twirapp/twir/apps/api-gql/internal/services/spotify_integration"
	spotify_song_requests "github.com/twirapp/twir/apps/api-gql/internal/services/spotify_song_requests"
	"github.com/twirapp/twir/apps/api-gql/internal/services/streamelements"
	streamlabsintegration "github.com/twirapp/twir/apps/api-gql/internal/services/streamlabs_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/toxic_messages"
	twir_events "github.com/twirapp/twir/apps/api-gql/internal/services/twir-events"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	valorantintegrationservice "github.com/twirapp/twir/apps/api-gql/internal/services/valorant_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	vkintegration "github.com/twirapp/twir/apps/api-gql/internal/services/vk_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/webhook_notifications"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	channelsongrequestssettingscache "github.com/twirapp/twir/libs/cache/channel_song_requests_settings"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	eventscache "github.com/twirapp/twir/libs/cache/channels_events_with_operations"
	channelsgamesvotebancache "github.com/twirapp/twir/libs/cache/channels_games_voteban"
	channelsintegrationssettingsseventvcache "github.com/twirapp/twir/libs/cache/channels_integrations_settings_seventv"
	channelsmoderationsettingsccahe "github.com/twirapp/twir/libs/cache/channels_moderation_settings"
	chattranslationssettignscache "github.com/twirapp/twir/libs/cache/chat_translations_settings"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	giveawayscache "github.com/twirapp/twir/libs/cache/giveaways"
	greetingscache "github.com/twirapp/twir/libs/cache/greetings"
	keywordscacher "github.com/twirapp/twir/libs/cache/keywords"
	quotescacher "github.com/twirapp/twir/libs/cache/quotes"
	rolescache "github.com/twirapp/twir/libs/cache/roles"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	valorantintegration "github.com/twirapp/twir/libs/integrations/valorant"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/wsrouter"

	channelsoverlaysservice "github.com/twirapp/twir/apps/api-gql/internal/services/channels_overlays"

	commandshttp "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/commands"

	twirhttp "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/twir"
)

const Service = "api-gql"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	OverlaysKappagenProviderSet,
	OverlaysBeRightBackProviderSet,
	OverlaysStreamStatsProviderSet,
	OverlaysTTSProviderSet,
	repositoriesSet,
	// services
	wire.NewSet(
		kickplatform.New,
		twitchplatform.New,
		vkvideo.NewBotSetupProvider,
		youtubeplatform.New,
		NewPlatformRegistry,
		channelservice.NewChannelService,
		NewValorantClient,
		dashboard_widget_events.New,
		dashboard_widgets.New,
		clientinfo.New,
		variables.New,
		timers.New,
		keywords.New,
		quotes.New,
		audit_logs.New,
		admin_actions.New,
		badges.New,
		badges_users.New,
		badges_with_users.New,
		users.New,
		twir_users.New,
		alerts.New,
		commands_with_groups_and_responses.New,
		commands_groups.New,
		commands_responses.New,
		commands.New,
		greetings.New,
		roles.New,
		roles_users.New,
		roles_with_roles_users.New,
		twitch.New,
		channels.New,
		channelplatformservice.NewFx,
		wire.Bind(new(resolvers.ChannelPlatformBindingsService), new(*channelplatformservice.Service)),
		dashboardaccess.NewFx,
		mcpOAuthService.NewFx,
		wire.Bind(new(mcpdelivery.AccessTokenVerifier), new(*mcpOAuthService.Service)),
		chat_messages.New,
		channels_commands_prefix.New,
		channels_emotes_usages.New,
		channels_secret.New,
		channels_storage.New,
		song_requests.New,
		spotify_song_requests.New,
		song_requests.NewPlaybackStateService,
		songrequestoverlaysettings.New,
		community_redemptions.New,
		streamelements.New,
		dashboard.New,
		seventv_integration.New,
		spotify_integration.New,
		scheduledvips.New,
		chat_wall.New,
		chat_translation.New,
		shortlinkscustomdomains.New,
		shortenedurls.New,
		uploaderservice.New,
		giveaways.New,
		overlays_dudes.New,
		channels_moderation_settings.New,
		pastebinsservice.New,
		events.New,
		twir_events.New,
		donatepay_integration.New,
		valorantintegrationservice.New,
		gamesvoteban.New,
		nightbotintegration.New,
		discord_integration.New,
		lastfmintegration.New,
		obs_websocket_module.New,
		webhook_notifications.New,
	),
	wire.NewSet(
		toxic_messages.New,
		channels_files.New,
		channels_redemptions_history.New,
		donationalertsintegration.New,
		donatestreamintegration.New,
		donatellointegration.New,
		vkintegration.New,
		faceitintegration.New,
		channelsoverlaysservice.New,
	),
	// app itself
	wire.NewSet(
		rate_limiter.NewLeakyBucket,
		httpmiddlewares.New,
		app.NewHuma,
		dataloader.New,
		auth.NewSessions,
		authroutes.New,
		minio.New,
		minio.NewUploaderS3,
		twitchcache.New,
		channelcache.New,
		channelscommandsprefixcache.New,
		greetingscache.New,
		commandscache.New,
		keywordscacher.New,
		quotescacher.New,
		giveawayscache.New,
		chatalertscache.New,
		channelalertscache.New,
		ttscache.NewTTSSettings,
		channelsmoderationsettingsccahe.New,
		chattranslationssettignscache.New,
		channelsongrequestssettingscache.New,
		channelsintegrationssettingsseventvcache.New,
		channelsgamesvotebancache.New,
		eventscache.New,
		rolescache.New,
		streamlabsintegration.New,
		wsrouter.NewNatsWsRouter,
		wire.Bind(new(wsrouter.WsRouter), new(*wsrouter.WsRouterNats)),
		twir_stats.New,
		resolvers.New,
		wire.Bind(new(resolvers.SelectedDashboardGetter), new(*auth.Auth)),
		wire.Bind(new(resolvers.CurrentPlatformGetter), new(*auth.Auth)),
		wire.Bind(new(resolvers.SessionReader), new(*auth.Auth)),
		directives.New,
		middlewares.New,
		server.New,
		mcpdelivery.New,
	),
	shortlinks.RegisterRoutes,
	uploaderoutes.RegisterRoutes,
	pastebins.RegisterRoutes,
	commandshttp.RegisterRoutes,
	ttsroutes.RegisterRoutes,
	brb.RegisterRoutes,
	twirhttp.RegisterRoutes,
	scheduledvipsroutes.RegisterRoutes,
	mcpOAuthRoutes.NewFromOpts,
	mcpOAuthRoutes.RegisterRoutes,
	gql.New,
	publicroutes.New,
	v2publicroutes.New,
	http_webhooks.New,
	song_requests.NewBridge,
	spotify_song_requests.NewReconciler,
	spotify_song_requests.NewRequestBridge,
	RegisterChannelsFilesRoute,
	RegisterValorantRoute,
	RegisterStreamRoute,
	RegisterMCP,
	wire.Struct(new(ApplicationDeps), "*"),
	NewApplication,
)

type ChannelsFilesRouteRegistration struct{}
type ValorantRouteRegistration struct{}
type StreamRouteRegistration struct{}
type MCPRegistration struct{}

func RegisterChannelsFilesRoute(
	api huma.API,
	config cfg.Config,
	channelsFilesService *channels_files.Service,
) ChannelsFilesRouteRegistration {
	channelsfilesroute.New(api, config, channelsFilesService)
	return ChannelsFilesRouteRegistration{}
}

func RegisterValorantRoute(
	api huma.API,
	config cfg.Config,
	sessions *auth.Auth,
	service *valorantintegrationservice.Service,
	kvClient kv.KV,
) ValorantRouteRegistration {
	valorant.New(api, config, sessions, service, kvClient)
	return ValorantRouteRegistration{}
}

func RegisterStreamRoute(
	streamsRepository streamsrepository.Repository,
	api huma.API,
	sessions *auth.Auth,
) StreamRouteRegistration {
	stream.New(streamsRepository, api, sessions)
	return StreamRouteRegistration{}
}

func RegisterMCP(s *server.Server, handler *mcpdelivery.Handler) MCPRegistration {
	mcpdelivery.Register(s, handler)
	return MCPRegistration{}
}

type ApplicationDeps struct {
	Lifecycle                     *lifecycle.Lifecycle
	PlatformRegistry              *platform.Registry
	GQL                           *gql.Gql
	PublicRoutes                  *publicroutes.Public
	V2PublicRoutes                *v2publicroutes.Public
	Webhooks                      *http_webhooks.Webhooks
	AuthRoutes                    *authroutes.Auth
	ChannelsFilesRoute            ChannelsFilesRouteRegistration
	SongRequestsBridge            *song_requests.Bridge
	SpotifySongRequestsReconciler *spotify_song_requests.Reconciler
	SpotifySongRequestsBridge     *spotify_song_requests.RequestBridge
	ValorantRoute                 ValorantRouteRegistration
	StreamRoute                   StreamRouteRegistration
	MCP                           MCPRegistration
	ShortlinksRoutes              shortlinks.Registration
	UploaderRoutes                uploaderoutes.Registration
	PastebinsRoutes               pastebins.Registration
	CommandsRoutes                commandshttp.Registration
	TTSRoutes                     ttsroutes.Registration
	BeRightBackRoutes             brb.Registration
	TwirRoutes                    twirhttp.Registration
	ScheduledVIPsRoutes           scheduledvipsroutes.Registration
	MCPOAuthRoutes                mcpOAuthRoutes.Registration
	Logger                        *slog.Logger
}

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(deps ApplicationDeps) *Application {
	deps.Logger.Info("🚀 API-GQL is running")
	return &Application{lifecycle: deps.Lifecycle}
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}

func NewValorantClient(config cfg.Config) *valorantintegration.HenrikValorantApiClient {
	return valorantintegration.NewHenrikApiClient(config.Valorant.HenrikApiKey)
}

func NewPlatformRegistry(
	config cfg.Config,
	twitchProvider *twitchplatform.Provider,
	kickProvider *kickplatform.Provider,
	youtubeProvider *youtubeplatform.Provider,
) (*platform.Registry, error) {
	return platform.NewFeatureGatedRegistry(
		config.IsVkVideoEnabled(),
		config.IsYouTubeEnabled(),
		[]platform.PlatformProvider{twitchProvider, kickProvider},
		func() (platform.PlatformProvider, error) {
			provider, err := vkvideo.New(vkvideo.Opts{Config: config})
			if err != nil {
				return nil, fmt.Errorf("create VK Video platform provider: %w", err)
			}

			return provider, nil
		},
		func() (platform.PlatformProvider, error) {
			return youtubeProvider, nil
		},
	)
}
