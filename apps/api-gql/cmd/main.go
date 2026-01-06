package main

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/app"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	publicroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public"
	http_webhooks "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-webhooks"
	httpmiddlewares "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	channelsfilesroute "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/channels/channels_files"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/integrations/valorant"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/overlays/brb"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/shortlinks"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/stream"
	ttsroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/tts"
	"github.com/twirapp/twir/apps/api-gql/internal/di"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_commands_prefix"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_emotes_usages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_files"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_moderation_settings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_redemptions_history"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_messages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_translation"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/community_redemptions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/dashboard"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
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
	nightbotintegration "github.com/twirapp/twir/apps/api-gql/internal/services/nightbot_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/obs_websocket_module"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays_dudes"
	pastebinsservice "github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_with_roles_users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/scheduled_vips"
	"github.com/twirapp/twir/apps/api-gql/internal/services/seventv_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/apps/api-gql/internal/services/song_requests"
	"github.com/twirapp/twir/apps/api-gql/internal/services/spotify_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/streamelements"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/toxic_messages"
	twir_events "github.com/twirapp/twir/apps/api-gql/internal/services/twir-events"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	valorantintegrationservice "github.com/twirapp/twir/apps/api-gql/internal/services/valorant_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	vkintegration "github.com/twirapp/twir/apps/api-gql/internal/services/vk_integration"
	"github.com/twirapp/twir/apps/parser/pkg/executron"
	"github.com/twirapp/twir/libs/baseapp"
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
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	channelsredemptionshistoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_redemptions_history/datasources/clickhouse"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
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
	overlaysdudesrepository "github.com/twirapp/twir/libs/repositories/overlays_dudes"
	overlaysdudesrepositorypgx "github.com/twirapp/twir/libs/repositories/overlays_dudes/pgx"
	rolesrepository "github.com/twirapp/twir/libs/repositories/roles"
	rolesrepositorypgx "github.com/twirapp/twir/libs/repositories/roles/pgx"
	rolesusersrepository "github.com/twirapp/twir/libs/repositories/roles_users"
	rolesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/roles_users/pgx"
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
	"github.com/twirapp/twir/libs/wsrouter"

	seventvintegrationrepository "github.com/twirapp/twir/libs/repositories/seventv_integration"
	seventvintegrationpostgres "github.com/twirapp/twir/libs/repositories/seventv_integration/datasource/postgres"

	botsrepository "github.com/twirapp/twir/libs/repositories/bots"
	botspostgres "github.com/twirapp/twir/libs/repositories/bots/datasource/postgres"

	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypostgres "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"

	integrationsrepository "github.com/twirapp/twir/libs/repositories/integrations"
	integrationspostgres "github.com/twirapp/twir/libs/repositories/integrations/datasource/postgres"

	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallpostgres "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"

	chattranslationrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	chattranslationpostgres "github.com/twirapp/twir/libs/repositories/chat_translation/datasource/postgres"

	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"

	channelsgiveawaysrepository "github.com/twirapp/twir/libs/repositories/giveaways"
	channelsgiveawaysrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways/pgx"

	channelsgiveawaysparticipantsrepository "github.com/twirapp/twir/libs/repositories/giveaways_participants"
	channelsgiveawaysparticipantsrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways_participants/pgx"

	channelsmoderationsettingsrepository "github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	channelsmoderationsettingsrepositorypostgres "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/datasource/postgres"

	pastebinsrepository "github.com/twirapp/twir/libs/repositories/pastebins"
	pastebinsrepositorypgx "github.com/twirapp/twir/libs/repositories/pastebins/datasource/postgres"

	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	toxicmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/toxic_messages/pgx"

	channelsfilesrepository "github.com/twirapp/twir/libs/repositories/channels_files"
	channelsfilesrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_files/datasource/postgres"

	plansrepository "github.com/twirapp/twir/libs/repositories/plans"
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

	"go.uber.org/fx"

	commandshttp "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/commands"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "api-gql",
			},
		),
		di.OverlaysKappagenModule,
		di.OverlaysBeRightBackModule,
		di.OverlaysTTSModule,
		// repositories
		fx.Provide(
			fx.Annotate(
				timersrepositorypgx.NewFx,
				fx.As(new(timersrepository.Repository)),
			),
			fx.Annotate(
				variablespgx.NewFx,
				fx.As(new(variablesrepository.Repository)),
			),
			fx.Annotate(
				keywordsrepositorypgx.NewFx,
				fx.As(new(keywordsrepository.Repository)),
			),
			fx.Annotate(
				channelsrepositorypgx.NewFx,
				fx.As(new(channelsrepository.Repository)),
			),
			fx.Annotate(
				badgesrepositorypgx.NewFx,
				fx.As(new(badgesrepository.Repository)),
			),
			fx.Annotate(
				badgesusersrepositorypgx.NewFx,
				fx.As(new(badgesusersrepository.Repository)),
			),
			fx.Annotate(
				usersrepositorypgx.NewFx,
				fx.As(new(usersrepository.Repository)),
			),
			fx.Annotate(
				userswithchannelrepositorypgx.NewFx,
				fx.As(new(userswithchannelrepository.Repository)),
			),
			fx.Annotate(
				alertsrepositorypgx.NewFx,
				fx.As(new(alertsrepository.Repository)),
			),
			fx.Annotate(
				commandswithgroupsandresponsesrepositorypgx.NewFx,
				fx.As(new(commandswithgroupsandresponsesrepository.Repository)),
			),
			fx.Annotate(
				commandsgroupsrepositorypgx.NewFx,
				fx.As(new(commandsgroupsrepository.Repository)),
			),
			fx.Annotate(
				commandsresponserepositorypgx.NewFx,
				fx.As(new(commandsresponserepository.Repository)),
			),
			fx.Annotate(
				commandsrepositorypgx.NewFx,
				fx.As(new(commandsrepository.Repository)),
			),
			fx.Annotate(
				rolesrepositorypgx.NewFx,
				fx.As(new(rolesrepository.Repository)),
			),
			fx.Annotate(
				rolesusersrepositorypgx.NewFx,
				fx.As(new(rolesusersrepository.Repository)),
			),
			fx.Annotate(
				greetingsrepositorypgx.NewFx,
				fx.As(new(greetingsrepository.Repository)),
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
				channelsintegrationsspotifypgx.NewFx,
				fx.As(new(channelsintegrationsspotify.Repository)),
			),
			fx.Annotate(
				seventvintegrationpostgres.NewFx,
				fx.As(new(seventvintegrationrepository.Repository)),
			),
			fx.Annotate(
				botspostgres.NewFx,
				fx.As(new(botsrepository.Repository)),
			),
			fx.Annotate(
				integrationspostgres.NewFx,
				fx.As(new(integrationsrepository.Repository)),
			),
			fx.Annotate(
				scheduledvipsrepositorypostgres.NewFx,
				fx.As(new(scheduledvipsrepository.Repository)),
			),
			fx.Annotate(
				chatwallpostgres.NewFx,
				fx.As(new(chatwallrepository.Repository)),
			),
			fx.Annotate(
				chattranslationpostgres.NewFx,
				fx.As(new(chattranslationrepository.Repository)),
			),
			fx.Annotate(
				shortenedurlsrepositorypostgres.NewFx,
				fx.As(new(shortenedurlsrepository.Repository)),
			),
			fx.Annotate(
				channelsgiveawaysparticipantsrepositorypgx.NewFx,
				fx.As(new(channelsgiveawaysparticipantsrepository.Repository)),
			),
			fx.Annotate(
				channelsgiveawaysrepositorypgx.NewFx,
				fx.As(new(channelsgiveawaysrepository.Repository)),
			),
			fx.Annotate(
				channelsmoderationsettingsrepositorypostgres.NewFx,
				fx.As(new(channelsmoderationsettingsrepository.Repository)),
			),
			fx.Annotate(
				overlaysdudesrepositorypgx.NewFx,
				fx.As(new(overlaysdudesrepository.Repository)),
			),
			fx.Annotate(
				eventsrepositorypgx.NewFx,
				fx.As(new(eventsrepository.Repository)),
			),
			fx.Annotate(
				pastebinsrepositorypgx.NewFx,
				fx.As(new(pastebinsrepository.Repository)),
			),
			fx.Annotate(
				toxicmessagesrepositorypgx.NewFx,
				fx.As(new(toxicmessagesrepository.Repository)),
			),
			fx.Annotate(
				channelsfilesrepositorypgx.NewFx,
				fx.As(new(channelsfilesrepository.Repository)),
			),
			fx.Annotate(
				plansrepositorypgx.NewFx,
				fx.As(new(plansrepository.Repository)),
			),
			fx.Annotate(
				channelsemotesusagesrepositoryclickhouse.NewFx,
				fx.As(new(channelsemotesusagesrepository.Repository)),
			),
			fx.Annotate(
				channelscommandsusagesclickhouse.NewFx,
				fx.As(new(channelscommandsusages.Repository)),
			),
			fx.Annotate(
				channelsredemptionshistoryclickhouse.NewFx,
				fx.As(new(channelsredemptionshistory.Repository)),
			),
			fx.Annotate(
				donatepayrepositorypostgres.NewFx,
				fx.As(new(donatepayrepository.Repository)),
			),
			fx.Annotate(
				tokensrepositorypgx.NewFx,
				fx.As(new(tokensrepository.Repository)),
			),
			fx.Annotate(
				channelsintegrationsvalorantpostgres.NewFx,
				fx.As(new(channelsintegrationsvalorant.Repository)),
			),
			fx.Annotate(
				streamsrepositorypostgres.NewFx,
				fx.As(new(streamsrepository.Repository)),
			),
			fx.Annotate(
				donationalertsrepoitorypostgres.NewFx,
				fx.As(new(donationalertsrepository.Repository)),
			),
			fx.Annotate(
				faceitrepositorypostgres.NewFx,
				fx.As(new(faceitrepository.Repository)),
			),
			fx.Annotate(
				channelsgamesvotebanpgx.NewFx,
				fx.As(new(channelsgamesvotebanrepository.Repository)),
			),
			fx.Annotate(
				channelsintegrationspostgres.NewFx,
				fx.As(new(channelsintegrationsrepository.Repository)),
			),
			fx.Annotate(
				channelsintegrationsdiscordpostgres.NewFx,
				fx.As(new(channelsintegrationsdiscordrepository.Repository)),
			),
			fx.Annotate(
				channelsintegrationslastfmpostgres.NewFx,
				fx.As(new(channelsintegrationslastfm.Repository)),
			),
			fx.Annotate(
				channelsmodulesobswebsocketpgx.NewFx,
				fx.As(new(channelsmodulesobswebsocket.Repository)),
			),
			fx.Annotate(
				vkintegrationrepopostgres.NewFx,
				fx.As(new(vkintegrationrepo.Repository)),
			),
			fx.Annotate(
				channelsoverlaysrepositorypgx.NewFx,
				fx.As(new(channelsoverlaysrepository.Repository)),
			),
		),
		// services
		fx.Provide(
			func(c cfg.Config) *valorantintegration.HenrikValorantApiClient {
				return valorantintegration.NewHenrikApiClient(c.Valorant.HenrikApiKey)
			},
			executron.New,
			dashboard_widget_events.New,
			variables.New,
			timers.New,
			keywords.New,
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
			chat_messages.New,
			channels_commands_prefix.New,
			channels_emotes_usages.New,
			song_requests.New,
			community_redemptions.New,
			streamelements.New,
			dashboard.New,
			seventv_integration.New,
			spotify_integration.New,
			scheduled_vips.New,
			chat_wall.New,
			chat_translation.New,
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
		),
		fx.Provide(
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
		fx.Provide(
			rate_limiter.NewLeakyBucket,
			httpmiddlewares.New,
			app.NewHuma,
			dataloader.New,
			auth.NewSessions,
			minio.New,
			twitchcache.New,
			channelcache.New,
			channelscommandsprefixcache.New,
			greetingscache.New,
			commandscache.New,
			keywordscacher.New,
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
			fx.Annotate(
				wsrouter.NewNatsWsRouterFx,
				fx.As(new(wsrouter.WsRouter)),
			),
			twir_stats.New,
			resolvers.New,
			directives.New,
			middlewares.New,
			server.New,
		),
		// huma routes
		shortlinks.FxModule,
		pastebins.FxModule,
		commandshttp.FxModule,
		ttsroutes.FxModule,
		brb.FxModule,
		// huma routes end
		fx.Invoke(
			gql.New,
			publicroutes.New,
			http_webhooks.New,
			httpbase.RegisterRoutes,
			authroutes.New,
			channelsfilesroute.New,
			valorant.New,
			stream.New,
			func(l *slog.Logger) {
				l.Info("🚀 API-GQL is running")
			},
		),
	).Run()
}
