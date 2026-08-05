package resolvers

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twirapp/kv"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	channelsservice "github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_commands_prefix"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_emotes_usages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_files"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_moderation_settings"
	channels_overlays "github.com/twirapp/twir/apps/api-gql/internal/services/channels_overlays"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_redemptions_history"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_secret"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_storage"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_messages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_translation"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
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
	"github.com/twirapp/twir/apps/api-gql/internal/services/donationalerts_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/events"
	faceitintegration "github.com/twirapp/twir/apps/api-gql/internal/services/faceit_integration"
	gamesvoteban "github.com/twirapp/twir/apps/api-gql/internal/services/games_voteban"
	"github.com/twirapp/twir/apps/api-gql/internal/services/giveaways"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	lastfmintegration "github.com/twirapp/twir/apps/api-gql/internal/services/lastfm_integration"
	nightbotintegration "github.com/twirapp/twir/apps/api-gql/internal/services/nightbot_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/obs_websocket_module"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/kappagen"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/tts"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays_dudes"
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
	"github.com/twirapp/twir/apps/api-gql/internal/services/streamlabs_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/toxic_messages"
	twir_events "github.com/twirapp/twir/apps/api-gql/internal/services/twir-events"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	twitchservice "github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	valorantintegration "github.com/twirapp/twir/apps/api-gql/internal/services/valorant_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	vkintegration "github.com/twirapp/twir/apps/api-gql/internal/services/vk_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/webhook_notifications"
	"github.com/twirapp/twir/libs/audit"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channelpublicsettingsrepo "github.com/twirapp/twir/libs/repositories/channel_public_settings"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channels_giveaways_settings "github.com/twirapp/twir/libs/repositories/channels_giveaways_settings"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	plansrepository "github.com/twirapp/twir/libs/repositories/plans"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	vkintegrationrepo "github.com/twirapp/twir/libs/repositories/vk_integration"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/wsrouter"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Deps struct {
	Logger        *slog.Logger
	AuditRecorder audit.Recorder
	WsRouter      wsrouter.WsRouter

	SpotifyRepository               channelsintegrationsspotify.Repository
	LastfmRepository                channelsintegrationslastfm.Repository
	VKIntegrationRepository         vkintegrationrepo.Repository
	PlansRepository                 plansrepository.Repository
	GiveawaysSettingsRepository     channels_giveaways_settings.Repository
	ChannelsRepository              channelsrepository.Repository
	UsersRepository                 usersrepository.Repository
	ChannelPublicSettingsRepository channelpublicsettingsrepo.Repository
	ChannelService                  *channelservice.ChannelService
	ChannelPlatformBindingsService  ChannelPlatformBindingsService
	ChannelPlatformDashboard        SelectedDashboardGetter
	CurrentPlatform                 CurrentPlatformGetter

	Sessions                         SessionReader
	Auth                             *authroutes.Auth
	Gorm                             *gorm.DB
	CachedTwitchClient               *twitchcahe.CachedTwitchClient
	CachedCommandsClient             *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	ChannelSongRequestsSettingsCache *generic_cacher.GenericCacher[model.ChannelSongRequestsSettings]
	Minio                            *minio.Client
	TwirBus                          *bus_core.Bus
	KV                               kv.KV
	TwirStats                        *twir_stats.TwirStats
	KickProvider                     *kickplatform.Provider

	DashboardWidgetEventsService          *dashboard_widget_events.Service
	DashboardWidgetsService               *dashboard_widgets.Service
	DashboardAccess                       *dashboardaccess.Service
	VariablesService                      *variables.Service
	TimersService                         *timers.Service
	KeywordsService                       *keywords.Service
	QuotesService                         *quotes.Service
	AuditLogsService                      *audit_logs.Service
	AdminActionsService                   *admin_actions.Service
	BadgesService                         *badges.Service
	BadgesUsersService                    *badges_users.Service
	UsersService                          *users.Service
	TwirUsersService                      *twir_users.Service
	AlertsService                         *alerts.Service
	CommandsService                       *commands.Service
	CommandsWithGroupsAndResponsesService *commands_with_groups_and_responses.Service
	CommandsResponsesService              *commands_responses.Service
	GreetingsService                      *greetings.Service
	RolesService                          *roles.Service
	RolesUsersService                     *roles_users.Service
	RolesWithUsersService                 *roles_with_roles_users.Service
	TwitchService                         *twitchservice.Service
	ChatMessagesService                   *chat_messages.Service
	ChannelsService                       *channelsservice.Service
	ChannelsCommandsPrefix                *channels_commands_prefix.Service
	ChannelsEmotesUsagesService           *channels_emotes_usages.Service
	TTSService                            *tts.Service
	SongRequestsService                   *song_requests.Service
	SongRequestPlaybackStateService       *song_requests.PlaybackStateService
	SongRequestOverlaySettingsService     *songrequestoverlaysettings.Service
	CommunityRedemptionsService           *community_redemptions.Service
	StreamElementsService                 *streamelements.Service
	DashboardService                      *dashboard.Service
	SevenTvIntegrationService             *seventv_integration.Service
	SpotifyIntegrationService             *spotify_integration.Service
	SpotifySongRequestsService            *spotify_song_requests.Service
	ValorantIntegrationService            *valorantintegration.Service
	DonatelloIntegrationService           *donatellointegration.Service
	DonateStreamIntegrationService        *donatestreamintegration.Service
	DiscordIntegrationService             *discord_integration.Service
	ScheduledVipsService                  *scheduledvips.Service
	ChatTranslationService                *chat_translation.Service
	ChatWallService                       *chat_wall.Service
	Config                                config.Config
	GiveawaysService                      *giveaways.Service
	ChannelsModerationSettingsService     *channels_moderation_settings.Service
	ShortenedUrlsService                  *shortenedurls.Service
	ShortLinksCustomDomainsService        *shortlinkscustomdomains.Service
	ToxicMessagesService                  *toxic_messages.Service
	ChatAlertsCache                       *generic_cacher.GenericCacher[chatalertscache.ChatAlert]
	ChannelsFilesService                  *channels_files.Service
	ChannelsRedemptionsHistoryService     *channels_redemptions_history.Service
	OverlaysDudesService                  *overlays_dudes.Service
	EventsService                         *events.Service
	KappagenService                       *kappagen.Service
	BeRightBackService                    *be_right_back.Service
	TwirEventsService                     *twir_events.Service
	DonatePayService                      *donatepay_integration.Service
	DonationAlertsIntegrationService      *donationalerts_integration.Service
	GamesVotebanService                   *gamesvoteban.Service
	NightbotIntegrationService            *nightbotintegration.Service
	LastfmIntegrationService              *lastfmintegration.Service
	ObsWebsocketModuleService             *obs_websocket_module.Service
	WebhookNotificationsService           *webhook_notifications.Service
	VKIntegrationService                  *vkintegration.Service
	FaceitIntegrationService              *faceitintegration.Service
	ChannelOverlaysService                *channels_overlays.Service
	StreamlabsIntegrationService          *streamlabs_integration.Service
	ChannelsSecretService                 *channels_secret.Service
	ChannelsStorageService                *channels_storage.Service
}

type ChannelPlatformBindingsService interface {
	List(context.Context, uuid.UUID) ([]channelplatformservice.Binding, error)
	Options() []channelplatformservice.Option
	Connect(context.Context, uuid.UUID, platformentity.Platform) (string, error)
	Disconnect(context.Context, uuid.UUID, platformentity.Platform) error
	SetEnabled(context.Context, uuid.UUID, platformentity.Platform, bool) (channelplatformservice.Binding, error)
}

type SelectedDashboardGetter interface {
	GetSelectedDashboard(context.Context) (string, error)
}

type CurrentPlatformGetter interface {
	GetCurrentPlatform(context.Context) (string, error)
}

type SessionReader interface {
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
	GetCurrentPlatform(context.Context) (string, error)
	GetSelectedDashboard(context.Context) (string, error)
	SetSessionSelectedDashboard(context.Context, string) error
	SessionLogout(context.Context) error
	GetChannelFromApiKey(context.Context) (channelentity.Channel, error)
}

type Resolver struct {
	deps Deps
}

func New(logger *slog.Logger, auditRecorder audit.Recorder, wsRouter wsrouter.WsRouter, spotifyRepository channelsintegrationsspotify.Repository, lastfmRepository channelsintegrationslastfm.Repository, vkIntegrationRepository vkintegrationrepo.Repository, plansRepository plansrepository.Repository, giveawaysSettingsRepository channels_giveaways_settings.Repository, channelsRepository channelsrepository.Repository, usersRepository usersrepository.Repository, channelPublicSettingsRepository channelpublicsettingsrepo.Repository, channelService *channelservice.ChannelService, channelPlatformBindingsService ChannelPlatformBindingsService, channelPlatformDashboard SelectedDashboardGetter, currentPlatform CurrentPlatformGetter, sessions SessionReader, authService *authroutes.Auth, gorm *gorm.DB, cachedTwitchClient *twitchcahe.CachedTwitchClient, cachedCommandsClient *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses], channelSongRequestsSettingsCache *generic_cacher.GenericCacher[model.ChannelSongRequestsSettings], minioClient *minio.Client, twirBus *bus_core.Bus, kv kv.KV, twirStats *twir_stats.TwirStats, kickProvider *kickplatform.Provider, dashboardWidgetEventsService *dashboard_widget_events.Service, dashboardWidgetsService *dashboard_widgets.Service, dashboardAccess *dashboardaccess.Service, variablesService *variables.Service, timersService *timers.Service, keywordsService *keywords.Service, quotesService *quotes.Service, auditLogsService *audit_logs.Service, adminActionsService *admin_actions.Service, badgesService *badges.Service, badgesUsersService *badges_users.Service, usersService *users.Service, twirUsersService *twir_users.Service, alertsService *alerts.Service, commandsService *commands.Service, commandsWithGroupsAndResponsesService *commands_with_groups_and_responses.Service, commandsResponsesService *commands_responses.Service, greetingsService *greetings.Service, rolesService *roles.Service, rolesUsersService *roles_users.Service, rolesWithUsersService *roles_with_roles_users.Service, twitchService *twitchservice.Service, chatMessagesService *chat_messages.Service, channelsService *channelsservice.Service, channelsCommandsPrefix *channels_commands_prefix.Service, channelsEmotesUsagesService *channels_emotes_usages.Service, ttsService *tts.Service, songRequestsService *song_requests.Service, songRequestPlaybackStateService *song_requests.PlaybackStateService, songRequestOverlaySettingsService *songrequestoverlaysettings.Service, communityRedemptionsService *community_redemptions.Service, streamElementsService *streamelements.Service, dashboardService *dashboard.Service, sevenTVIntegrationService *seventv_integration.Service, spotifyIntegrationService *spotify_integration.Service, spotifySongRequestsService *spotify_song_requests.Service, valorantIntegrationService *valorantintegration.Service, donatelloIntegrationService *donatellointegration.Service, donateStreamIntegrationService *donatestreamintegration.Service, discordIntegrationService *discord_integration.Service, scheduledVipsService *scheduledvips.Service, chatTranslationService *chat_translation.Service, chatWallService *chat_wall.Service, cfg config.Config, giveawaysService *giveaways.Service, channelsModerationSettingsService *channels_moderation_settings.Service, shortenedUrlsService *shortenedurls.Service, shortLinksCustomDomainsService *shortlinkscustomdomains.Service, toxicMessagesService *toxic_messages.Service, chatAlertsCache *generic_cacher.GenericCacher[chatalertscache.ChatAlert], channelsFilesService *channels_files.Service, channelsRedemptionsHistoryService *channels_redemptions_history.Service, overlaysDudesService *overlays_dudes.Service, eventsService *events.Service, kappagenService *kappagen.Service, beRightBackService *be_right_back.Service, twirEventsService *twir_events.Service, donatePayService *donatepay_integration.Service, donationAlertsIntegrationService *donationalerts_integration.Service, gamesVotebanService *gamesvoteban.Service, nightbotIntegrationService *nightbotintegration.Service, lastfmIntegrationService *lastfmintegration.Service, obsWebsocketModuleService *obs_websocket_module.Service, webhookNotificationsService *webhook_notifications.Service, vkIntegrationService *vkintegration.Service, faceitIntegrationService *faceitintegration.Service, channelOverlaysService *channels_overlays.Service, streamlabsIntegrationService *streamlabs_integration.Service, channelsSecretService *channels_secret.Service, channelsStorageService *channels_storage.Service) (*Resolver, error) {
	return &Resolver{
		deps: Deps{
			Logger: logger, AuditRecorder: auditRecorder, WsRouter: wsRouter,
			SpotifyRepository: spotifyRepository, LastfmRepository: lastfmRepository,
			VKIntegrationRepository: vkIntegrationRepository, PlansRepository: plansRepository,
			GiveawaysSettingsRepository: giveawaysSettingsRepository, ChannelsRepository: channelsRepository,
			UsersRepository: usersRepository, ChannelPublicSettingsRepository: channelPublicSettingsRepository,
			ChannelService: channelService, ChannelPlatformBindingsService: channelPlatformBindingsService,
			ChannelPlatformDashboard: channelPlatformDashboard, CurrentPlatform: currentPlatform,
			Sessions: sessions, Auth: authService, Gorm: gorm, CachedTwitchClient: cachedTwitchClient,
			CachedCommandsClient: cachedCommandsClient, ChannelSongRequestsSettingsCache: channelSongRequestsSettingsCache,
			Minio: minioClient, TwirBus: twirBus, KV: kv, TwirStats: twirStats, KickProvider: kickProvider,
			DashboardWidgetEventsService: dashboardWidgetEventsService, DashboardWidgetsService: dashboardWidgetsService,
			DashboardAccess: dashboardAccess, VariablesService: variablesService, TimersService: timersService,
			KeywordsService: keywordsService, QuotesService: quotesService, AuditLogsService: auditLogsService,
			AdminActionsService: adminActionsService, BadgesService: badgesService, BadgesUsersService: badgesUsersService,
			UsersService: usersService, TwirUsersService: twirUsersService, AlertsService: alertsService,
			CommandsService: commandsService, CommandsWithGroupsAndResponsesService: commandsWithGroupsAndResponsesService,
			CommandsResponsesService: commandsResponsesService, GreetingsService: greetingsService, RolesService: rolesService,
			RolesUsersService: rolesUsersService, RolesWithUsersService: rolesWithUsersService, TwitchService: twitchService,
			ChatMessagesService: chatMessagesService, ChannelsService: channelsService,
			ChannelsCommandsPrefix: channelsCommandsPrefix, ChannelsEmotesUsagesService: channelsEmotesUsagesService,
			TTSService: ttsService, SongRequestsService: songRequestsService,
			SongRequestPlaybackStateService:   songRequestPlaybackStateService,
			SongRequestOverlaySettingsService: songRequestOverlaySettingsService,
			CommunityRedemptionsService:       communityRedemptionsService, StreamElementsService: streamElementsService,
			DashboardService: dashboardService, SevenTvIntegrationService: sevenTVIntegrationService,
			SpotifyIntegrationService: spotifyIntegrationService, SpotifySongRequestsService: spotifySongRequestsService,
			ValorantIntegrationService:     valorantIntegrationService,
			DonatelloIntegrationService:    donatelloIntegrationService,
			DonateStreamIntegrationService: donateStreamIntegrationService, DiscordIntegrationService: discordIntegrationService,
			ScheduledVipsService: scheduledVipsService, ChatTranslationService: chatTranslationService,
			ChatWallService: chatWallService, Config: cfg, GiveawaysService: giveawaysService,
			ChannelsModerationSettingsService: channelsModerationSettingsService,
			ShortenedUrlsService:              shortenedUrlsService, ShortLinksCustomDomainsService: shortLinksCustomDomainsService,
			ToxicMessagesService: toxicMessagesService, ChatAlertsCache: chatAlertsCache,
			ChannelsFilesService: channelsFilesService, ChannelsRedemptionsHistoryService: channelsRedemptionsHistoryService,
			OverlaysDudesService: overlaysDudesService, EventsService: eventsService, KappagenService: kappagenService,
			BeRightBackService: beRightBackService, TwirEventsService: twirEventsService, DonatePayService: donatePayService,
			DonationAlertsIntegrationService: donationAlertsIntegrationService, GamesVotebanService: gamesVotebanService,
			NightbotIntegrationService: nightbotIntegrationService, LastfmIntegrationService: lastfmIntegrationService,
			ObsWebsocketModuleService: obsWebsocketModuleService, WebhookNotificationsService: webhookNotificationsService,
			VKIntegrationService: vkIntegrationService, FaceitIntegrationService: faceitIntegrationService,
			ChannelOverlaysService: channelOverlaysService, StreamlabsIntegrationService: streamlabsIntegrationService,
			ChannelsSecretService: channelsSecretService, ChannelsStorageService: channelsStorageService,
		},
	}, nil
}

func NewFromDeps(deps Deps) (*Resolver, error) {
	return &Resolver{deps: deps}, nil
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(
	ctx *graphql.OperationContext,
	fields []graphql.CollectedField,
	prefix string,
) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(
			preloads,
			GetNestedPreloads(
				ctx,
				graphql.CollectFields(ctx, column.Selections, nil),
				prefixColumn,
			)...,
		)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
