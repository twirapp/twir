package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/goforj/wire"
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
	mcpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/di"
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
	"github.com/twirapp/twir/apps/api-gql/internal/services/streamelements"
	streamlabsintegration "github.com/twirapp/twir/apps/api-gql/internal/services/streamlabs_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/toxic_messages"
	twir_events "github.com/twirapp/twir/apps/api-gql/internal/services/twir-events"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
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
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	badgesrepository "github.com/twirapp/twir/libs/repositories/badges"
	badgesrepositorypgx "github.com/twirapp/twir/libs/repositories/badges/pgx"
	badgesusersrepository "github.com/twirapp/twir/libs/repositories/badges_users"
	badgesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/badges_users/pgx"
	channelplatformsrepository "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsrepositorypgx "github.com/twirapp/twir/libs/repositories/channel_platforms/pgx"
	channelpublicsettingsrepo "github.com/twirapp/twir/libs/repositories/channel_public_settings"
	channelpublicsettingspgx "github.com/twirapp/twir/libs/repositories/channel_public_settings/datasource/postgres"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsemotesusagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/datasources/clickhouse"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	channelsintegrationslastfmpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm/datasources/postgres"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	channelsintegrationsvalorant "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant"
	channelsintegrationsvalorantpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/datasources/postgres"
	channelsmodulesobswebsocket "github.com/twirapp/twir/libs/repositories/channels_modules_obs_websocket"
	channelsmodulesobswebsocketpgx "github.com/twirapp/twir/libs/repositories/channels_modules_obs_websocket/datasources/postgres"
	channelsmoduleswebhooks "github.com/twirapp/twir/libs/repositories/channels_modules_webhooks"
	channelsmoduleswebhookspgx "github.com/twirapp/twir/libs/repositories/channels_modules_webhooks/pgx"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	channelsredemptionshistoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_redemptions_history/datasources/clickhouse"
	channelssecretrepository "github.com/twirapp/twir/libs/repositories/channels_secret"
	channelssecretpgx "github.com/twirapp/twir/libs/repositories/channels_secret/pgx"
	channelsstoragerepository "github.com/twirapp/twir/libs/repositories/channels_storage"
	channelsstoragepgx "github.com/twirapp/twir/libs/repositories/channels_storage/pgx"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
	"github.com/twirapp/twir/libs/repositories/command_role_cooldown"
	commandrolecooldownpgx "github.com/twirapp/twir/libs/repositories/command_role_cooldown/pgx"
	commandsrepository "github.com/twirapp/twir/libs/repositories/commands"
	commandsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands/pgx"
	commandsgroupsrepository "github.com/twirapp/twir/libs/repositories/commands_group"
	commandsgroupsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands_group/pgx"
	commandsresponserepository "github.com/twirapp/twir/libs/repositories/commands_response"
	commandsresponserepositorypgx "github.com/twirapp/twir/libs/repositories/commands_response/pgx"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsesrepositorypgx "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
	keywordsrepositorypgx "github.com/twirapp/twir/libs/repositories/keywords/pgx"
	mcpOAuthRepository "github.com/twirapp/twir/libs/repositories/mcp_oauth"
	mcpOAuthRepositoryPostgres "github.com/twirapp/twir/libs/repositories/mcp_oauth/datasource/postgres"
	overlaysdudesrepository "github.com/twirapp/twir/libs/repositories/overlays_dudes"
	overlaysdudesrepositorypgx "github.com/twirapp/twir/libs/repositories/overlays_dudes/pgx"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
	quotesrepositorypgx "github.com/twirapp/twir/libs/repositories/quotes/pgx"
	rolesrepository "github.com/twirapp/twir/libs/repositories/roles"
	rolesrepositorypgx "github.com/twirapp/twir/libs/repositories/roles/pgx"
	rolesusersrepository "github.com/twirapp/twir/libs/repositories/roles_users"
	rolesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/roles_users/pgx"
	shortlinksbanneduapresetpatternsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_preset_patterns"
	shortlinksbanneduapresetpatternsrepositorypgx "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_preset_patterns/datasource/postgres"
	shortlinksbanneduapresetsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_presets"
	shortlinksbanneduapresetsrepositorypgx "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_presets/datasource/postgres"
	shortlinkscustomdomainsrepository "github.com/twirapp/twir/libs/repositories/short_links_custom_domains"
	shortlinkscustomdomainsrepositorypgx "github.com/twirapp/twir/libs/repositories/short_links_custom_domains/pgx"
	shortlinkslinkbannedusaragentsrepository "github.com/twirapp/twir/libs/repositories/short_links_link_banned_user_agents"
	shortlinkslinkbannedusaragentsrepositorypgx "github.com/twirapp/twir/libs/repositories/short_links_link_banned_user_agents/datasource/postgres"
	shortlinkslinkpresetsrepository "github.com/twirapp/twir/libs/repositories/short_links_link_presets"
	shortlinkslinkpresetsrepositorypgx "github.com/twirapp/twir/libs/repositories/short_links_link_presets/datasource/postgres"
	shortlinksviewsrepository "github.com/twirapp/twir/libs/repositories/short_links_views"
	shortlinksviewsrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/short_links_views/datasources/clickhouse"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	shortenedurlsrepositorypostgres "github.com/twirapp/twir/libs/repositories/shortened_urls/datasource/postgres"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	userswithchannelrepository "github.com/twirapp/twir/libs/repositories/users_with_channel"
	userswithchannelrepositorypgx "github.com/twirapp/twir/libs/repositories/users_with_channel/pgx"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	variablespgx "github.com/twirapp/twir/libs/repositories/variables/pgx"
	vkintegrationrepo "github.com/twirapp/twir/libs/repositories/vk_integration"
	vkintegrationrepopostgres "github.com/twirapp/twir/libs/repositories/vk_integration/datasource/postgres"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/wsrouter"

	seventvintegrationrepository "github.com/twirapp/twir/libs/repositories/seventv_integration"
	seventvintegrationpostgres "github.com/twirapp/twir/libs/repositories/seventv_integration/datasource/postgres"

	botsrepository "github.com/twirapp/twir/libs/repositories/bots"
	botspostgres "github.com/twirapp/twir/libs/repositories/bots/datasource/postgres"

	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypostgres "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"

	integrationsrepository "github.com/twirapp/twir/libs/repositories/integrations"
	integrationspostgres "github.com/twirapp/twir/libs/repositories/integrations/datasource/postgres"

	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	vkvideobotsrepository "github.com/twirapp/twir/libs/repositories/vk_video_bots"
	vkvideobotsrepositorypgx "github.com/twirapp/twir/libs/repositories/vk_video_bots/datasource/postgres"
	youtubebotsrepository "github.com/twirapp/twir/libs/repositories/youtube_bots"
	youtubebotsrepositorypgx "github.com/twirapp/twir/libs/repositories/youtube_bots/datasource/postgres"

	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallpostgres "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"

	chattranslationrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	chattranslationpostgres "github.com/twirapp/twir/libs/repositories/chat_translation/datasource/postgres"

	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	songrequestoverlaysettingsrepository "github.com/twirapp/twir/libs/repositories/song_request_overlay_settings"
	songrequestoverlaysettingspgx "github.com/twirapp/twir/libs/repositories/song_request_overlay_settings/pgx"

	channelsgiveawaysrepository "github.com/twirapp/twir/libs/repositories/giveaways"
	channelsgiveawaysrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways/pgx"

	channelsgiveawaysparticipantsrepository "github.com/twirapp/twir/libs/repositories/giveaways_participants"
	channelsgiveawaysparticipantsrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways_participants/pgx"

	channelsgiveawayssettingsrepository "github.com/twirapp/twir/libs/repositories/channels_giveaways_settings"
	channelsgiveawayssettingsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_giveaways_settings/pgx"

	channelsmoderationsettingsrepository "github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	channelsmoderationsettingsrepositorypostgres "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/datasource/postgres"

	pastebinsrepository "github.com/twirapp/twir/libs/repositories/pastebins"
	pastebinsrepositorypgx "github.com/twirapp/twir/libs/repositories/pastebins/datasource/postgres"

	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	toxicmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/toxic_messages/pgx"

	channelsfilesrepository "github.com/twirapp/twir/libs/repositories/channels_files"
	channelsfilesrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_files/datasource/postgres"

	plansrepositorypgx "github.com/twirapp/twir/libs/repositories/plans/pgx"

	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
	channelscommandsusagesclickhouse "github.com/twirapp/twir/libs/repositories/channels_commands_usages/datasources/clickhouse"

	eventsrepository "github.com/twirapp/twir/libs/repositories/events"
	eventsrepositorypgx "github.com/twirapp/twir/libs/repositories/events/pgx"

	donatepayrepository "github.com/twirapp/twir/libs/repositories/donatepay_integration"
	donatepayrepositorypostgres "github.com/twirapp/twir/libs/repositories/donatepay_integration/datasource/postgres"

	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokensrepositorypgx "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"

	donationalertsrepository "github.com/twirapp/twir/libs/repositories/donationalerts_integration"
	donationalertsrepoitorypostgres "github.com/twirapp/twir/libs/repositories/donationalerts_integration/datasource/postgres"
	faceitrepository "github.com/twirapp/twir/libs/repositories/faceit_integration"
	faceitrepositorypostgres "github.com/twirapp/twir/libs/repositories/faceit_integration/datasource/postgres"

	channelsoverlaysservice "github.com/twirapp/twir/apps/api-gql/internal/services/channels_overlays"
	channelsoverlaysrepository "github.com/twirapp/twir/libs/repositories/channels_overlays"
	channelsoverlaysrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_overlays/pgx"

	channelsgamesvotebanrepository "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	channelsgamesvotebanpgx "github.com/twirapp/twir/libs/repositories/channels_games_voteban/pgx"

	channelsintegrationsrepository "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationspostgres "github.com/twirapp/twir/libs/repositories/channels_integrations/datasource/postgres"

	channelsintegrationsdiscordrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	channelsintegrationsdiscordpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_discord/datasource/postgres"

	commandshttp "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/commands"

	twirhttp "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/twir"

	streamlabsrepository "github.com/twirapp/twir/libs/repositories/streamlabs_integration"
	streamlabsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streamlabs_integration/datasource/postgres"

	dashboardwidgetsrepository "github.com/twirapp/twir/libs/repositories/dashboard_widgets"
	dashboardwidgetsrepositorypgx "github.com/twirapp/twir/libs/repositories/dashboard_widgets/pgx"
)

const Service = "api-gql"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	parameterSet,
	di.OverlaysKappagenProviderSet,
	di.OverlaysBeRightBackProviderSet,
	di.OverlaysTTSProviderSet,
	// repositories
	wire.NewSet(
		timersrepositorypgx.NewFx,
		wire.Bind(new(timersrepository.Repository), new(*timersrepositorypgx.Pgx)),
		variablespgx.NewFx,
		wire.Bind(new(variablesrepository.Repository), new(*variablespgx.Pgx)),
		channelssecretpgx.NewFx,
		wire.Bind(new(channelssecretrepository.Repository), new(*channelssecretpgx.Pgx)),
		channelsstoragepgx.NewFx,
		wire.Bind(new(channelsstoragerepository.Repository), new(*channelsstoragepgx.Pgx)),
		keywordsrepositorypgx.NewFx,
		wire.Bind(new(keywordsrepository.Repository), new(*keywordsrepositorypgx.Pgx)),
		quotesrepositorypgx.NewFx,
		wire.Bind(new(quotesrepository.Repository), new(*quotesrepositorypgx.Pgx)),
		channelsrepositorypgx.NewFx,
		wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
		channelplatformsrepositorypgx.NewFx,
		wire.Bind(new(channelplatformsrepository.Repository), new(*channelplatformsrepositorypgx.Pgx)),
		channelpublicsettingspgx.NewFx,
		wire.Bind(new(channelpublicsettingsrepo.Repository), new(*channelpublicsettingspgx.Pgx)),
		badgesrepositorypgx.NewFx,
		wire.Bind(new(badgesrepository.Repository), new(*badgesrepositorypgx.Pgx)),
		badgesusersrepositorypgx.NewFx,
		wire.Bind(new(badgesusersrepository.Repository), new(*badgesusersrepositorypgx.Pgx)),
		usersrepositorypgx.NewFx,
		wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
		mcpOAuthRepositoryPostgres.NewFx,
		wire.Bind(new(mcpOAuthRepository.Repository), new(*mcpOAuthRepositoryPostgres.Pgx)),
		userswithchannelrepositorypgx.NewFx,
		wire.Bind(new(userswithchannelrepository.Repository), new(*userswithchannelrepositorypgx.Pgx)),
		alertsrepositorypgx.NewFx,
		wire.Bind(new(alertsrepository.Repository), new(*alertsrepositorypgx.Pgx)),
		commandswithgroupsandresponsesrepositorypgx.NewFx,
		wire.Bind(new(commandswithgroupsandresponsesrepository.Repository), new(*commandswithgroupsandresponsesrepositorypgx.Pgx)),
		commandsgroupsrepositorypgx.NewFx,
		wire.Bind(new(commandsgroupsrepository.Repository), new(*commandsgroupsrepositorypgx.Pgx)),
		commandsresponserepositorypgx.NewFx,
		wire.Bind(new(commandsresponserepository.Repository), new(*commandsresponserepositorypgx.Pgx)),
		commandsrepositorypgx.NewFx,
		wire.Bind(new(commandsrepository.Repository), new(*commandsrepositorypgx.Pgx)),
		rolesrepositorypgx.NewFx,
		wire.Bind(new(rolesrepository.Repository), new(*rolesrepositorypgx.Pgx)),
		rolesusersrepositorypgx.NewFx,
		wire.Bind(new(rolesusersrepository.Repository), new(*rolesusersrepositorypgx.Pgx)),
		greetingsrepositorypgx.NewFx,
		wire.Bind(new(greetingsrepository.Repository), new(*greetingsrepositorypgx.Pgx)),
		chatmessagesrepositoryclickhouse.NewFx,
		wire.Bind(new(chatmessagesrepository.Repository), new(*chatmessagesrepositoryclickhouse.Clickhouse)),
		channelscommandsprefixpgx.NewFx,
		wire.Bind(new(channelscommandsprefixrepository.Repository), new(*channelscommandsprefixpgx.Pgx)),
		channelsintegrationsspotifypgx.NewFx,
		wire.Bind(new(channelsintegrationsspotify.Repository), new(*channelsintegrationsspotifypgx.Pgx)),
		seventvintegrationpostgres.NewFx,
		wire.Bind(new(seventvintegrationrepository.Repository), new(*seventvintegrationpostgres.Pgx)),
		botspostgres.NewFx,
		wire.Bind(new(botsrepository.Repository), new(*botspostgres.Pgx)),
		integrationspostgres.NewFx,
		wire.Bind(new(integrationsrepository.Repository), new(*integrationspostgres.Pgx)),
		kickbotsrepositorypgx.NewFx,
		wire.Bind(new(kickbotsrepository.Repository), new(*kickbotsrepositorypgx.Pgx)),
		vkvideobotsrepositorypgx.NewFx,
		wire.Bind(new(vkvideobotsrepository.Repository), new(*vkvideobotsrepositorypgx.Pgx)),
		youtubebotsrepositorypgx.NewFx,
		wire.Bind(new(youtubebotsrepository.Repository), new(*youtubebotsrepositorypgx.Pgx)),
		chatwallpostgres.NewFx,
		wire.Bind(new(chatwallrepository.Repository), new(*chatwallpostgres.Pgx)),
		scheduledvipsrepositorypostgres.NewFx,
		wire.Bind(new(scheduledvipsrepository.Repository), new(*scheduledvipsrepositorypostgres.Pgx)),
		chattranslationpostgres.NewFx,
		wire.Bind(new(chattranslationrepository.Repository), new(*chattranslationpostgres.Pgx)),
		shortenedurlsrepositorypostgres.NewFx,
		wire.Bind(new(shortenedurlsrepository.Repository), new(*shortenedurlsrepositorypostgres.Pgx)),
		shortlinkscustomdomainsrepositorypgx.NewFx,
		wire.Bind(new(shortlinkscustomdomainsrepository.Repository), new(*shortlinkscustomdomainsrepositorypgx.Pgx)),
		shortlinksbanneduapresetsrepositorypgx.NewFx,
		wire.Bind(new(shortlinksbanneduapresetsrepository.Repository), new(*shortlinksbanneduapresetsrepositorypgx.Pgx)),
		shortlinksbanneduapresetpatternsrepositorypgx.NewFx,
		wire.Bind(new(shortlinksbanneduapresetpatternsrepository.Repository), new(*shortlinksbanneduapresetpatternsrepositorypgx.Pgx)),
		shortlinkslinkpresetsrepositorypgx.NewFx,
		wire.Bind(new(shortlinkslinkpresetsrepository.Repository), new(*shortlinkslinkpresetsrepositorypgx.Pgx)),
		shortlinkslinkbannedusaragentsrepositorypgx.NewFx,
		wire.Bind(new(shortlinkslinkbannedusaragentsrepository.Repository), new(*shortlinkslinkbannedusaragentsrepositorypgx.Pgx)),
		channelsgiveawaysparticipantsrepositorypgx.NewFx,
		wire.Bind(new(channelsgiveawaysparticipantsrepository.Repository), new(*channelsgiveawaysparticipantsrepositorypgx.Pgx)),
		channelsgiveawaysrepositorypgx.NewFx,
		wire.Bind(new(channelsgiveawaysrepository.Repository), new(*channelsgiveawaysrepositorypgx.Pgx)),
		channelsgiveawayssettingsrepositorypgx.NewFx,
		wire.Bind(new(channelsgiveawayssettingsrepository.Repository), new(*channelsgiveawayssettingsrepositorypgx.Pgx)),
		channelsmoderationsettingsrepositorypostgres.NewFx,
		wire.Bind(new(channelsmoderationsettingsrepository.Repository), new(*channelsmoderationsettingsrepositorypostgres.Pgx)),
		overlaysdudesrepositorypgx.NewFx,
		wire.Bind(new(overlaysdudesrepository.Repository), new(*overlaysdudesrepositorypgx.Pgx)),
		eventsrepositorypgx.NewFx,
		wire.Bind(new(eventsrepository.Repository), new(*eventsrepositorypgx.Pgx)),
		pastebinsrepositorypgx.NewFx,
		wire.Bind(new(pastebinsrepository.Repository), new(*pastebinsrepositorypgx.Pgx)),
		toxicmessagesrepositorypgx.NewFx,
		wire.Bind(new(toxicmessagesrepository.Repository), new(*toxicmessagesrepositorypgx.Pgx)),
		channelsfilesrepositorypgx.NewFx,
		wire.Bind(new(channelsfilesrepository.Repository), new(*channelsfilesrepositorypgx.Pgx)),
		plansrepositorypgx.NewFx,
		dashboardwidgetsrepositorypgx.NewFx,
		wire.Bind(new(dashboardwidgetsrepository.Repository), new(*dashboardwidgetsrepositorypgx.Pgx)),
		channelsemotesusagesrepositoryclickhouse.NewFx,
		wire.Bind(new(channelsemotesusagesrepository.Repository), new(*channelsemotesusagesrepositoryclickhouse.Clickhouse)),
		channelscommandsusagesclickhouse.NewFx,
		wire.Bind(new(channelscommandsusages.Repository), new(*channelscommandsusagesclickhouse.Clickhouse)),
		channelsredemptionshistoryclickhouse.NewFx,
		wire.Bind(new(channelsredemptionshistory.Repository), new(*channelsredemptionshistoryclickhouse.Clickhouse)),
		shortlinksviewsrepositoryclickhouse.NewFx,
		wire.Bind(new(shortlinksviewsrepository.Repository), new(*shortlinksviewsrepositoryclickhouse.Clickhouse)),
		donatepayrepositorypostgres.NewFx,
		wire.Bind(new(donatepayrepository.Repository), new(*donatepayrepositorypostgres.Pgx)),
		tokensrepositorypgx.NewFx,
		wire.Bind(new(tokensrepository.Repository), new(*tokensrepositorypgx.Pgx)),
		channelsintegrationsvalorantpostgres.NewFx,
		wire.Bind(new(channelsintegrationsvalorant.Repository), new(*channelsintegrationsvalorantpostgres.Pgx)),
		streamsrepositorypostgres.NewFx,
		wire.Bind(new(streamsrepository.Repository), new(*streamsrepositorypostgres.Pgx)),
		donationalertsrepoitorypostgres.NewFx,
		wire.Bind(new(donationalertsrepository.Repository), new(*donationalertsrepoitorypostgres.Pgx)),
		faceitrepositorypostgres.NewFx,
		wire.Bind(new(faceitrepository.Repository), new(*faceitrepositorypostgres.Pgx)),
		channelsgamesvotebanpgx.NewFx,
		wire.Bind(new(channelsgamesvotebanrepository.Repository), new(*channelsgamesvotebanpgx.Pgx)),
		channelsintegrationspostgres.NewFx,
		wire.Bind(new(channelsintegrationsrepository.Repository), new(*channelsintegrationspostgres.Pgx)),
		channelsintegrationsdiscordpostgres.NewFx,
		wire.Bind(new(channelsintegrationsdiscordrepository.Repository), new(*channelsintegrationsdiscordpostgres.Pgx)),
		channelsintegrationslastfmpostgres.NewFx,
		wire.Bind(new(channelsintegrationslastfm.Repository), new(*channelsintegrationslastfmpostgres.Pgx)),
		channelsmodulesobswebsocketpgx.NewFx,
		wire.Bind(new(channelsmodulesobswebsocket.Repository), new(*channelsmodulesobswebsocketpgx.Pgx)),
		channelsmoduleswebhookspgx.NewFx,
		wire.Bind(new(channelsmoduleswebhooks.Repository), new(*channelsmoduleswebhookspgx.Pgx)),
		vkintegrationrepopostgres.NewFx,
		wire.Bind(new(vkintegrationrepo.Repository), new(*vkintegrationrepopostgres.Pgx)),
		channelsoverlaysrepositorypgx.NewFx,
		wire.Bind(new(channelsoverlaysrepository.Repository), new(*channelsoverlaysrepositorypgx.Pgx)),
		streamlabsrepositorypostgres.NewFx,
		wire.Bind(new(streamlabsrepository.Repository), new(*streamlabsrepositorypostgres.Pgx)),
		commandrolecooldownpgx.NewFx,
		wire.Bind(new(command_role_cooldown.Repository), new(*commandrolecooldownpgx.Pgx)),
		songrequestoverlaysettingspgx.NewFx,
		wire.Bind(new(songrequestoverlaysettingsrepository.Repository), new(*songrequestoverlaysettingspgx.Pgx)),
	),
	// services
	wire.NewSet(
		kickplatform.New,
		twitchplatform.New,
		vkvideo.NewBotSetupProvider,
		youtubeplatform.New,
		newPlatformRegistry,
		channelservice.NewChannelService,
		newValorantClient,
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
	wire.Struct(new(shortlinks.Dependencies), "*"),
	shortlinks.RegisterRoutes,
	wire.Struct(new(pastebins.Dependencies), "*"),
	pastebins.RegisterRoutes,
	wire.Struct(new(commandshttp.Dependencies), "*"),
	commandshttp.RegisterRoutes,
	wire.Struct(new(ttsroutes.Dependencies), "*"),
	ttsroutes.RegisterRoutes,
	wire.Struct(new(brb.Dependencies), "*"),
	brb.RegisterRoutes,
	wire.Struct(new(twirhttp.Dependencies), "*"),
	twirhttp.RegisterRoutes,
	wire.Struct(new(scheduledvipsroutes.Dependencies), "*"),
	scheduledvipsroutes.RegisterRoutes,
	wire.Struct(new(mcpOAuthRoutes.ProviderOpts), "*"),
	mcpOAuthRoutes.NewFromOpts,
	mcpOAuthRoutes.RegisterRoutes,
	gql.New,
	publicroutes.New,
	v2publicroutes.New,
	http_webhooks.New,
	song_requests.NewBridge,
	registerChannelsFilesRoute,
	registerValorantRoute,
	registerStreamRoute,
	registerMCP,
	wire.Struct(new(ApplicationDeps), "*"),
	NewApplication,
)

var parameterSet = wire.NewSet(
	wire.Struct(new(app.HumaOpts), "*"),
	wire.Struct(new(auth.Opts), "*"),
	wire.Struct(new(dataloader.Opts), "*"),
	wire.Struct(new(directives.Opts), "*"),
	wire.Struct(new(gql.Opts), "*"),
	wire.Struct(new(resolvers.Deps), "*"),
	wire.Struct(new(twir_stats.Opts), "*"),
	wire.Struct(new(publicroutes.Opts), "*"),
	wire.Struct(new(v2publicroutes.Opts), "*"),
	wire.Struct(new(http_webhooks.Opts), "*"),
	wire.Struct(new(httpmiddlewares.Opts), "*"),
	wire.Struct(new(authroutes.Opts), "*"),
	wire.Struct(new(channelsfilesroute.Opts), "*"),
	wire.Struct(new(valorant.Opts), "*"),
	wire.Struct(new(stream.Opts), "*"),
	wire.Struct(new(mcpdelivery.Deps), "*"),
	wire.Struct(new(kickplatform.Opts), "*"),
	wire.Struct(new(twitchplatform.Opts), "*"),
	wire.Struct(new(vkvideo.BotSetupProviderOpts), "*"),
	wire.Struct(new(youtubeplatform.Opts), "*"),
	wire.Struct(new(middlewares.Opts), "*"),
	wire.Struct(new(server.Opts), "*"),
	wire.Struct(new(dashboard_widget_events.DashboardWidgetsEventsOpts), "*"),
	wire.Struct(new(dashboard_widgets.Opts), "*"),
	wire.Struct(new(clientinfo.Opts), "*"),
	wire.Struct(new(variables.Opts), "*"),
	wire.Struct(new(timers.Opts), "*"),
	wire.Struct(new(keywords.Opts), "*"),
	wire.Struct(new(quotes.Opts), "*"),
	wire.Struct(new(audit_logs.Opts), "*"),
	wire.Struct(new(admin_actions.Opts), "*"),
	wire.Struct(new(badges.Opts), "*"),
	wire.Struct(new(badges_users.Opts), "*"),
	wire.Struct(new(badges_with_users.Opts), "*"),
	wire.Struct(new(users.Opts), "*"),
	wire.Struct(new(twir_users.Opts), "*"),
	wire.Struct(new(alerts.Opts), "*"),
	wire.Struct(new(commands_with_groups_and_responses.Opts), "*"),
	wire.Struct(new(commands_groups.Opts), "*"),
	wire.Struct(new(commands_responses.Opts), "*"),
	wire.Struct(new(commands.Opts), "*"),
	wire.Struct(new(greetings.Opts), "*"),
	wire.Struct(new(roles.Opts), "*"),
	wire.Struct(new(roles_users.Opts), "*"),
	wire.Struct(new(roles_with_roles_users.Opts), "*"),
	wire.Struct(new(twitch.Opts), "*"),
	wire.Struct(new(channels.Opts), "*"),
	wire.Struct(new(channelplatformservice.Opts), "*"),
	wire.Struct(new(dashboardaccess.Opts), "*"),
	wire.Struct(new(mcpOAuthService.Opts), "*"),
	wire.Struct(new(chat_messages.Opts), "*"),
	wire.Struct(new(channels_commands_prefix.Opts), "*"),
	wire.Struct(new(channels_emotes_usages.Opts), "*"),
	wire.Struct(new(channels_secret.Opts), "*"),
	wire.Struct(new(channels_storage.Opts), "*"),
	wire.Struct(new(song_requests.Opts), "*"),
	wire.Struct(new(song_requests.PlaybackStateOpts), "*"),
	wire.Struct(new(song_requests.BridgeOpts), "*"),
	wire.Struct(new(songrequestoverlaysettings.Opts), "*"),
	wire.Struct(new(community_redemptions.Opts), "*"),
	wire.Struct(new(streamelements.Opts), "*"),
	wire.Struct(new(dashboard.Opts), "*"),
	wire.Struct(new(seventv_integration.Opts), "*"),
	wire.Struct(new(scheduledvips.Opts), "*"),
	wire.Struct(new(chat_wall.Opts), "*"),
	wire.Struct(new(chat_translation.Opts), "*"),
	wire.Struct(new(shortlinkscustomdomains.Opts), "*"),
	wire.Struct(new(shortenedurls.Opts), "*"),
	wire.Struct(new(giveaways.Opts), "*"),
	wire.Struct(new(overlays_dudes.Opts), "*"),
	wire.Struct(new(channels_moderation_settings.Opts), "*"),
	wire.Struct(new(pastebinsservice.Opts), "*"),
	wire.Struct(new(events.Opts), "*"),
	wire.Struct(new(twir_events.Opts), "*"),
	wire.Struct(new(donatepay_integration.Opts), "*"),
	wire.Struct(new(valorantintegrationservice.Opts), "*"),
	wire.Struct(new(gamesvoteban.Opts), "*"),
	wire.Struct(new(nightbotintegration.Opts), "*"),
	wire.Struct(new(discord_integration.Opts), "*"),
	wire.Struct(new(lastfmintegration.Opts), "*"),
	wire.Struct(new(obs_websocket_module.Opts), "*"),
	wire.Struct(new(webhook_notifications.Opts), "*"),
	wire.Struct(new(toxic_messages.Opts), "*"),
	wire.Struct(new(channels_files.Opts), "*"),
	wire.Struct(new(channels_redemptions_history.Opts), "*"),
	wire.Struct(new(donationalertsintegration.Opts), "*"),
	wire.Struct(new(donatestreamintegration.Opts), "*"),
	wire.Struct(new(donatellointegration.Opts), "*"),
	wire.Struct(new(vkintegration.Opts), "*"),
	wire.Struct(new(faceitintegration.Opts), "*"),
	wire.Struct(new(channelsoverlaysservice.Opts), "*"),
	wire.Struct(new(streamlabsintegration.Opts), "*"),
)

type ChannelsFilesRouteRegistration struct{}
type ValorantRouteRegistration struct{}
type StreamRouteRegistration struct{}
type MCPRegistration struct{}

func registerChannelsFilesRoute(opts channelsfilesroute.Opts) ChannelsFilesRouteRegistration {
	channelsfilesroute.New(opts)
	return ChannelsFilesRouteRegistration{}
}

func registerValorantRoute(opts valorant.Opts) ValorantRouteRegistration {
	valorant.New(opts)
	return ValorantRouteRegistration{}
}

func registerStreamRoute(opts stream.Opts) StreamRouteRegistration {
	stream.New(opts)
	return StreamRouteRegistration{}
}

func registerMCP(s *server.Server, handler *mcpdelivery.Handler) MCPRegistration {
	mcpdelivery.Register(s, handler)
	return MCPRegistration{}
}

type ApplicationDeps struct {
	Lifecycle           *lifecycle.Lifecycle
	PlatformRegistry    *platform.Registry
	GQL                 *gql.Gql
	PublicRoutes        *publicroutes.Public
	V2PublicRoutes      *v2publicroutes.Public
	Webhooks            *http_webhooks.Webhooks
	AuthRoutes          *authroutes.Auth
	ChannelsFilesRoute  ChannelsFilesRouteRegistration
	SongRequestsBridge  *song_requests.Bridge
	ValorantRoute       ValorantRouteRegistration
	StreamRoute         StreamRouteRegistration
	MCP                 MCPRegistration
	ShortlinksRoutes    shortlinks.Registration
	PastebinsRoutes     pastebins.Registration
	CommandsRoutes      commandshttp.Registration
	TTSRoutes           ttsroutes.Registration
	BeRightBackRoutes   brb.Registration
	TwirRoutes          twirhttp.Registration
	ScheduledVIPsRoutes scheduledvipsroutes.Registration
	MCPOAuthRoutes      mcpOAuthRoutes.Registration
	Logger              *slog.Logger
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

func main() {
	application, err := initializeApplication()
	if err != nil {
		log.Fatalf("initialize application: %v", err)
	}
	if err := application.Run(); err != nil {
		log.Fatalf("run application: %v", err)
	}
}

func newValorantClient(config cfg.Config) *valorantintegration.HenrikValorantApiClient {
	return valorantintegration.NewHenrikApiClient(config.Valorant.HenrikApiKey)
}

func newPlatformRegistry(
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
