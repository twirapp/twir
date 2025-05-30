package main

import (
	"github.com/satont/twir/apps/parser/pkg/executron"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/app"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	publicroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public"
	http_webhooks "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-webhooks"
	httpmiddlewares "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/shortlinks"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_commands_prefix"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_moderation_settings"
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
	"github.com/twirapp/twir/apps/api-gql/internal/services/giveaways"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/tts"
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
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	"github.com/twirapp/twir/libs/baseapp"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	channelsongrequestssettingscache "github.com/twirapp/twir/libs/cache/channel_song_requests_settings"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	channelsintegrationssettingsseventvcache "github.com/twirapp/twir/libs/cache/channels_integrations_settings_seventv"
	channelsmoderationsettingsccahe "github.com/twirapp/twir/libs/cache/channels_moderation_settings"
	chattranslationssettignscache "github.com/twirapp/twir/libs/cache/chat_translations_settings"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	giveawayscache "github.com/twirapp/twir/libs/cache/giveaways"
	greetingscache "github.com/twirapp/twir/libs/cache/greetings"
	keywordscacher "github.com/twirapp/twir/libs/cache/keywords"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	badgesrepository "github.com/twirapp/twir/libs/repositories/badges"
	badgesrepositorypgx "github.com/twirapp/twir/libs/repositories/badges/pgx"
	badgesusersrepository "github.com/twirapp/twir/libs/repositories/badges_users"
	badgesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/badges_users/pgx"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/chat_messages/pgx"
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
	rolesrepository "github.com/twirapp/twir/libs/repositories/roles"
	rolesrepositorypgx "github.com/twirapp/twir/libs/repositories/roles/pgx"
	rolesusersrepository "github.com/twirapp/twir/libs/repositories/roles_users"
	rolesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/roles_users/pgx"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	shortenedurlsrepositorypostgres "github.com/twirapp/twir/libs/repositories/shortened_urls/datasource/postgres"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	userswithchannelrepository "github.com/twirapp/twir/libs/repositories/users_with_channel"
	userswithchannelrepositorypgx "github.com/twirapp/twir/libs/repositories/users_with_channel/pgx"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	variablespgx "github.com/twirapp/twir/libs/repositories/variables/pgx"

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

	"go.uber.org/fx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "api-gql",
			},
		),
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
				chatmessagesrepositorypgx.NewFx,
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
				pastebinsrepositorypgx.NewFx,
				fx.As(new(pastebinsrepository.Repository)),
			),
			fx.Annotate(
				toxicmessagesrepositorypgx.NewFx,
				fx.As(new(toxicmessagesrepository.Repository)),
			),
		),
		// services
		fx.Provide(
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
			tts.New,
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
			channels_moderation_settings.New,
			pastebinsservice.New,
			toxic_messages.New,
		),
		// grpc clients
		fx.Provide(
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
		),
		// app itself
		fx.Provide(
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
			channelsmoderationsettingsccahe.New,
			chattranslationssettignscache.New,
			channelsongrequestssettingscache.New,
			channelsintegrationssettingsseventvcache.New,
			fx.Annotate(
				wsrouter.NewNatsSubscription,
				fx.As(new(wsrouter.WsRouter)),
			),
			twir_stats.New,
			resolvers.New,
			directives.New,
			middlewares.New,
			server.New,
		),
		fx.Invoke(
			gql.New,
			publicroutes.New,
			http_webhooks.New,
			authroutes.New,
			shortlinks.New,
			pastebins.New,
		),
	).Run()
}
