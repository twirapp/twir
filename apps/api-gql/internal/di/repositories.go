package di

import (
	"github.com/goforj/wire"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	badgesrepository "github.com/twirapp/twir/libs/repositories/badges"
	badgesrepositorypgx "github.com/twirapp/twir/libs/repositories/badges/pgx"
	badgesusersrepository "github.com/twirapp/twir/libs/repositories/badges_users"
	badgesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/badges_users/pgx"
	botsrepository "github.com/twirapp/twir/libs/repositories/bots"
	botspostgres "github.com/twirapp/twir/libs/repositories/bots/datasource/postgres"
	channelplatformsrepository "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsrepositorypgx "github.com/twirapp/twir/libs/repositories/channel_platforms/pgx"
	channelpublicsettingsrepo "github.com/twirapp/twir/libs/repositories/channel_public_settings"
	channelpublicsettingspgx "github.com/twirapp/twir/libs/repositories/channel_public_settings/datasource/postgres"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
	channelscommandsusagesclickhouse "github.com/twirapp/twir/libs/repositories/channels_commands_usages/datasources/clickhouse"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsemotesusagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/datasources/clickhouse"
	channelsfilesrepository "github.com/twirapp/twir/libs/repositories/channels_files"
	channelsfilesrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_files/datasource/postgres"
	channelsgamesvotebanrepository "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	channelsgamesvotebanpgx "github.com/twirapp/twir/libs/repositories/channels_games_voteban/pgx"
	channelsgiveawayssettingsrepository "github.com/twirapp/twir/libs/repositories/channels_giveaways_settings"
	channelsgiveawayssettingsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_giveaways_settings/pgx"
	channelsintegrationsrepository "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationspostgres "github.com/twirapp/twir/libs/repositories/channels_integrations/datasource/postgres"
	channelsintegrationsdiscordrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	channelsintegrationsdiscordpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_discord/datasource/postgres"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	channelsintegrationslastfmpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm/datasources/postgres"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	channelsintegrationsvalorant "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant"
	channelsintegrationsvalorantpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/datasources/postgres"
	channelsmoderationsettingsrepository "github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	channelsmoderationsettingsrepositorypostgres "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/datasource/postgres"
	channelsmodulesobswebsocket "github.com/twirapp/twir/libs/repositories/channels_modules_obs_websocket"
	channelsmodulesobswebsocketpgx "github.com/twirapp/twir/libs/repositories/channels_modules_obs_websocket/datasources/postgres"
	channelsmoduleswebhooks "github.com/twirapp/twir/libs/repositories/channels_modules_webhooks"
	channelsmoduleswebhookspgx "github.com/twirapp/twir/libs/repositories/channels_modules_webhooks/pgx"
	channelsoverlaysrepository "github.com/twirapp/twir/libs/repositories/channels_overlays"
	channelsoverlaysrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_overlays/pgx"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	channelsredemptionshistoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_redemptions_history/datasources/clickhouse"
	channelssecretrepository "github.com/twirapp/twir/libs/repositories/channels_secret"
	channelssecretpgx "github.com/twirapp/twir/libs/repositories/channels_secret/pgx"
	channelsstoragerepository "github.com/twirapp/twir/libs/repositories/channels_storage"
	channelsstoragepgx "github.com/twirapp/twir/libs/repositories/channels_storage/pgx"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
	chattranslationrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	chattranslationpostgres "github.com/twirapp/twir/libs/repositories/chat_translation/datasource/postgres"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallpostgres "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"
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
	dashboardwidgetsrepository "github.com/twirapp/twir/libs/repositories/dashboard_widgets"
	dashboardwidgetsrepositorypgx "github.com/twirapp/twir/libs/repositories/dashboard_widgets/pgx"
	donatepayrepository "github.com/twirapp/twir/libs/repositories/donatepay_integration"
	donatepayrepositorypostgres "github.com/twirapp/twir/libs/repositories/donatepay_integration/datasource/postgres"
	donationalertsrepository "github.com/twirapp/twir/libs/repositories/donationalerts_integration"
	donationalertsrepoitorypostgres "github.com/twirapp/twir/libs/repositories/donationalerts_integration/datasource/postgres"
	eventsrepository "github.com/twirapp/twir/libs/repositories/events"
	eventsrepositorypgx "github.com/twirapp/twir/libs/repositories/events/pgx"
	faceitrepository "github.com/twirapp/twir/libs/repositories/faceit_integration"
	faceitrepositorypostgres "github.com/twirapp/twir/libs/repositories/faceit_integration/datasource/postgres"
	channelsgiveawaysrepository "github.com/twirapp/twir/libs/repositories/giveaways"
	channelsgiveawaysrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways/pgx"
	channelsgiveawaysparticipantsrepository "github.com/twirapp/twir/libs/repositories/giveaways_participants"
	channelsgiveawaysparticipantsrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways_participants/pgx"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	integrationsrepository "github.com/twirapp/twir/libs/repositories/integrations"
	integrationspostgres "github.com/twirapp/twir/libs/repositories/integrations/datasource/postgres"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
	keywordsrepositorypgx "github.com/twirapp/twir/libs/repositories/keywords/pgx"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	mcpOAuthRepository "github.com/twirapp/twir/libs/repositories/mcp_oauth"
	mcpOAuthRepositoryPostgres "github.com/twirapp/twir/libs/repositories/mcp_oauth/datasource/postgres"
	overlaysdudesrepository "github.com/twirapp/twir/libs/repositories/overlays_dudes"
	overlaysdudesrepositorypgx "github.com/twirapp/twir/libs/repositories/overlays_dudes/pgx"
	pastebinsrepository "github.com/twirapp/twir/libs/repositories/pastebins"
	pastebinsrepositorypgx "github.com/twirapp/twir/libs/repositories/pastebins/datasource/postgres"
	plansrepositorypgx "github.com/twirapp/twir/libs/repositories/plans/pgx"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
	quotesrepositorypgx "github.com/twirapp/twir/libs/repositories/quotes/pgx"
	requestedsongsrepository "github.com/twirapp/twir/libs/repositories/requested_songs"
	requestedsongspostgres "github.com/twirapp/twir/libs/repositories/requested_songs/datasource/postgres"
	rolesrepository "github.com/twirapp/twir/libs/repositories/roles"
	rolesrepositorypgx "github.com/twirapp/twir/libs/repositories/roles/pgx"
	rolesusersrepository "github.com/twirapp/twir/libs/repositories/roles_users"
	rolesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/roles_users/pgx"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypostgres "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	seventvintegrationrepository "github.com/twirapp/twir/libs/repositories/seventv_integration"
	seventvintegrationpostgres "github.com/twirapp/twir/libs/repositories/seventv_integration/datasource/postgres"
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
	songrequestoverlaysettingsrepository "github.com/twirapp/twir/libs/repositories/song_request_overlay_settings"
	songrequestoverlaysettingspgx "github.com/twirapp/twir/libs/repositories/song_request_overlay_settings/pgx"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
	songrequestssettingspostgres "github.com/twirapp/twir/libs/repositories/song_requests_settings/datasource/postgres"
	spotifysongrequestsrepository "github.com/twirapp/twir/libs/repositories/spotify_song_requests"
	spotifysongrequestsrepositorypgx "github.com/twirapp/twir/libs/repositories/spotify_song_requests/datasource/postgres"
	streamlabsrepository "github.com/twirapp/twir/libs/repositories/streamlabs_integration"
	streamlabsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streamlabs_integration/datasource/postgres"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokensrepositorypgx "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"
	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	toxicmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/toxic_messages/pgx"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	userswithchannelrepository "github.com/twirapp/twir/libs/repositories/users_with_channel"
	userswithchannelrepositorypgx "github.com/twirapp/twir/libs/repositories/users_with_channel/pgx"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	variablespgx "github.com/twirapp/twir/libs/repositories/variables/pgx"
	vkintegrationrepo "github.com/twirapp/twir/libs/repositories/vk_integration"
	vkintegrationrepopostgres "github.com/twirapp/twir/libs/repositories/vk_integration/datasource/postgres"
	vkvideobotsrepository "github.com/twirapp/twir/libs/repositories/vk_video_bots"
	vkvideobotsrepositorypgx "github.com/twirapp/twir/libs/repositories/vk_video_bots/datasource/postgres"
	youtubebotsrepository "github.com/twirapp/twir/libs/repositories/youtube_bots"
	youtubebotsrepositorypgx "github.com/twirapp/twir/libs/repositories/youtube_bots/datasource/postgres"
)

// repositoriesSet wires all libs/repositories implementations to their interfaces.
var repositoriesSet = wire.NewSet(
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
	spotifysongrequestsrepositorypgx.NewFx,
	wire.Bind(new(spotifysongrequestsrepository.Repository), new(*spotifysongrequestsrepositorypgx.Pgx)),
	requestedsongspostgres.NewFx,
	wire.Bind(new(requestedsongsrepository.Repository), new(*requestedsongspostgres.Pgx)),
	songrequestssettingspostgres.NewFx,
	wire.Bind(new(songrequestssettingsrepository.Repository), new(*songrequestssettingspostgres.Pgx)),
)
