-- +goose Up
-- +goose StatementBegin

create extension if not exists "uuid-ossp";

SELECT 'up SQL query';
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'bots_type_enum') THEN create type bots_type_enum as enum ('DEFAULT', 'CUSTOM');
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_commands_cooldowntype_enum') THEN create type channels_commands_cooldowntype_enum as enum ('GLOBAL', 'PER_USER');
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_moderation_settings_type_enum') THEN create type channels_moderation_settings_type_enum as enum (
			'links', 'blacklists', 'symbols',
			'longMessage', 'caps', 'emotes'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'notifications_messages_langcode_enum') THEN create type notifications_messages_langcode_enum as enum ('RU', 'GB');
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channel_events_list_type_enum') THEN create type channel_events_list_type_enum as enum (
			'follow', 'subscription', 'resubscription',
			'donation', 'host', 'raid', 'moderator_added',
			'moderator_remove'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_customvars_type_enum') THEN create type channels_customvars_type_enum as enum ('SCRIPT', 'TEXT', 'NUMBER');
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_roles_type_enum') THEN create type channels_roles_type_enum as enum (
			'BROADCASTER', 'MODERATOR', 'SUBSCRIBER',
			'VIP', 'CUSTOM'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_roles_permissions_enum') THEN create type channels_roles_permissions_enum as enum (
			'CAN_ACCESS_DASHBOARD', 'UPDATE_CHANNEL_TITLE',
			'UPDATE_CHANNEL_CATEGORY', 'VIEW_COMMANDS',
			'MANAGE_COMMANDS', 'VIEW_KEYWORDS',
			'MANAGE_KEYWORDS', 'VIEW_TIMERS',
			'MANAGE_TIMERS', 'VIEW_INTEGRATIONS',
			'MANAGE_INTEGRATIONS', 'VIEW_SONG_REQUESTS',
			'MANAGE_SONG_REQUESTS', 'VIEW_MODERATION',
			'MANAGE_MODERATION', 'VIEW_VARIABLES',
			'MANAGE_VARIABLES', 'VIEW_GREETINGS',
			'MANAGE_GREETINGS'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_modules_settings_type_enum') THEN create type channels_modules_settings_type_enum as enum (
			'youtube_song_requests', 'obs_websocket',
			'tts'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_events_operations_filters_type_enum') THEN create type channels_events_operations_filters_type_enum as enum (
			'EQUALS', 'NOT_EQUALS', 'CONTAINS',
			'NOT_CONTAINS', 'STARTS_WITH', 'ENDS_WITH',
			'GREATER_THAN', 'LESS_THAN', 'GREATER_THAN_OR_EQUALS',
			'LESS_THAN_OR_EQUALS', 'IS_EMPTY',
			'IS_NOT_EMPTY'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_commands_module_enum') THEN create type channels_commands_module_enum as enum (
			'CUSTOM', 'DOTA', 'MODERATION', 'MANAGE',
			'SONGS', 'TTS', 'STATS'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_events_operations_type_enum') THEN create type channels_events_operations_type_enum as enum (
			'TIMEOUT', 'TIMEOUT_RANDOM', 'BAN',
			'UNBAN', 'BAN_RANDOM', 'VIP', 'UNVIP',
			'UNVIP_RANDOM', 'UNVIP_RANDOM_IF_NO_SLOTS',
			'MOD', 'UNMOD', 'UNMOD_RANDOM', 'SEND_MESSAGE',
			'CHANGE_TITLE', 'CHANGE_CATEGORY',
			'FULFILL_REDEMPTION', 'CANCEL_REDEMPTION',
			'ENABLE_SUBMODE', 'DISABLE_SUBMODE',
			'ENABLE_EMOTEONLY', 'DISABLE_EMOTEONLY',
			'CREATE_GREETING', 'OBS_SET_SCENE',
			'OBS_TOGGLE_SOURCE', 'OBS_TOGGLE_AUDIO',
			'OBS_AUDIO_SET_VOLUME', 'OBS_AUDIO_INCREASE_VOLUME',
			'OBS_AUDIO_DECREASE_VOLUME', 'OBS_DISABLE_AUDIO',
			'OBS_ENABLE_AUDIO', 'OBS_START_STREAM',
			'OBS_STOP_STREAM', 'CHANGE_VARIABLE',
			'INCREMENT_VARIABLE', 'DECREMENT_VARIABLE',
			'TTS_SAY', 'TTS_SKIP', 'TTS_ENABLE',
			'TTS_DISABLE', 'TTS_ENABLE_AUTOREAD',
			'TTS_DISABLE_AUTOREAD', 'TTS_SWITCH_AUTOREAD',
			'ALLOW_COMMAND_TO_USER', 'REMOVE_ALLOW_COMMAND_TO_USER',
			'DENY_COMMAND_TO_USER', 'REMOVE_DENY_COMMAND_TO_USER'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'channels_events_type_enum') THEN create type channels_events_type_enum as enum (
			'FOLLOW', 'SUBSCRIBE', 'RESUBSCRIBE',
			'SUB_GIFT', 'REDEMPTION_CREATED',
			'COMMAND_USED', 'FIRST_USER_MESSAGE',
			'RAIDED', 'TITLE_OR_CATEGORY_CHANGED',
			'STREAM_ONLINE', 'STREAM_OFFLINE',
			'ON_CHAT_CLEAR', 'DONATE', 'KEYWORD_MATCHED',
			'GREETING_SENDED', 'POLL_BEGIN',
			'POLL_PROGRESS', 'POLL_END', 'PREDICTION_BEGIN',
			'PREDICTION_PROGRESS', 'PREDICTION_END',
			'PREDICTION_LOCK', 'STREAM_FIRST_USER_JOIN'
			);
		END IF;
	END
$$;
DO
$$
	BEGIN
		IF NOT EXISTS (SELECT 1
									 FROM pg_type
									 WHERE typname = 'integrations_service_enum') THEN create type integrations_service_enum as enum (
			'LASTFM', 'VK', 'FACEIT', 'SPOTIFY',
			'DONATIONALERTS', 'STREAMLABS',
			'DONATEPAY', 'DONATELLO', 'VALORANT',
			'DONATE_STREAM'
			);
		END IF;
	END
$$;
create table if not exists dota_game_modes
(
	id   integer not null
		constraint "PK_d4eb4f935e0316f80207af11f49" primary key,
	name text    not null
);
create table if not exists dota_heroes
(
	id   integer not null
		constraint "PK_8f920265eebb91aa3ae6b0cb6e2" primary key,
	name text    not null
);
create index if not exists "IDX_82490264788f0b786e2ff312a5" on dota_heroes (name);
create table if not exists dota_matches
(
	id                            text    default gen_random_uuid() not null
		constraint "PK_279584095dc95716f433cc184f3" primary key,
	"startedAt"                   timestamp                         not null,
	lobby_type                    integer,
	players                       integer[],
	players_heroes                integer[],
	weekend_tourney_bracket_round text,
	weekend_tourney_skill_level   text,
	match_id                      text                              not null
		constraint "UQ_27d39ee5a1cdb04ed78a52286a1" unique,
	avarage_mmr                   integer                           not null,
	"lobbyId"                     text                              not null,
	finished                      boolean default false             not null,
	"gameModeId"                  integer                           not null
		constraint "FK_7fa32afcc462ae6f50b176f6c9c" references dota_game_modes on update cascade on delete restrict
);
create index if not exists "IDX_27d39ee5a1cdb04ed78a52286a" on dota_matches (match_id);
create table if not exists dota_matches_cards
(
	id               text default gen_random_uuid() not null
		constraint "PK_49fc741bc298a5e8402480bb0c9" primary key,
	account_id       text                           not null,
	rank_tier        integer,
	leaderboard_rank integer,
	match_id         text                           not null
		constraint "FK_a9d58830e97acdda53124a0431f" references dota_matches on update cascade on delete restrict
);
create unique index if not exists dota_matches_cards_match_id_account_id_key on dota_matches_cards (account_id, match_id);
create table if not exists dota_matches_results
(
	id          text default gen_random_uuid() not null
		constraint "PK_2c4f099e78034e57cef9a42b9fc" primary key,
	players     jsonb                          not null,
	radiant_win boolean                        not null,
	game_mode   integer                        not null,
	match_id    text                           not null
		constraint "REL_66717d38c93b88a283ad77736a" unique
		constraint "FK_66717d38c93b88a283ad77736ad" references dota_matches (match_id) on update cascade on delete restrict
);
create unique index if not exists dota_matches_results_match_id_key on dota_matches_results (match_id);
create table if not exists integrations
(
	id             text default gen_random_uuid() not null
		constraint "PK_9adcdc6d6f3922535361ce641e8" primary key,
	service        integrations_service_enum      not null,
	"accessToken"  text,
	"refreshToken" text,
	"clientId"     text,
	"clientSecret" text,
	"apiKey"       text,
	"redirectUrl"  text
);
create table if not exists tokens
(
	id                    uuid   default uuid_generate_v4() not null
		constraint "PK_3001e89ada36263dabf1fb6210a" primary key,
	"accessToken"         text                              not null,
	"refreshToken"        text                              not null,
	"expiresIn"           integer                           not null,
	"obtainmentTimestamp" timestamp with time zone          not null,
	scopes                text[] default '{}' :: text[]
);
create table if not exists bots
(
	id        text           not null primary key,
	type      bots_type_enum not null,
	"tokenId" uuid
		constraint "UQ_df4240a5d71aa6a23b829d3cee8" unique
		constraint "bots_tokenId_key" references tokens on
			update
			cascade on delete
			set
			null
);
create index if not exists "IDX_df4240a5d71aa6a23b829d3cee" on bots ("tokenId");
create table if not exists users
(
	id           text                              not null
		constraint "PK_a3ffb1c0c8416b9fc6f907b7433" primary key,
	"isTester"   boolean default false             not null,
	"isBotAdmin" boolean default false             not null,
	"tokenId"    uuid
		constraint "REL_d98a275f8bc6cd986fcbe2eab0" unique
		constraint "FK_d98a275f8bc6cd986fcbe2eab01" references tokens on
			update
			cascade on delete
			set
			null,
	"apiKey"     uuid    default gen_random_uuid() not null
);
create table if not exists channels
(
	id               text                  not null
		constraint "PK_bc603823f3f741359c2339389f9" primary key
		constraint "FK_bc603823f3f741359c2339389f9" references users on
			update
			cascade on delete restrict,
	"isEnabled"      boolean default true  not null,
	"isTwitchBanned" boolean default false not null,
	"isBanned"       boolean default false not null,
	"botId"          text                  not null
		constraint "FK_4f890144c0cb55fe7867b8f61e6" references bots on
			update
			cascade on delete restrict,
	"isBotMod"       boolean default false not null
);
create table if not exists channels_customvars
(
	id          uuid default uuid_generate_v4() not null
		constraint "PK_a6d0167a26c310ef7832c66b55d" primary key,
	name        text                            not null,
	description text,
	type        channels_customvars_type_enum   not null,
	"evalValue" text default '' :: text         not null,
	response    text default '' :: text         not null,
	"channelId" text                            not null
		constraint "FK_2bf1744a6cd76e4457eddfd3bc4" references channels on update cascade on delete restrict
);
create table if not exists channels_dota_accounts
(
	id          text not null,
	"channelId" text not null
		constraint "FK_a4a2cd666dac0cae74549a5de72" references channels on
			update
			cascade on delete restrict,
	constraint "PK_f325fe7c08cb3c3c6ccaaa509fa" primary key (id, "channelId")
);
create unique index if not exists "channels_dota_accounts_id_channelId_key" on channels_dota_accounts ("channelId", id);
create table if not exists channels_greetings
(
	id          uuid    default uuid_generate_v4() not null
		constraint "PK_90ee111c7fefdaaff4918cab6df" primary key,
	"userId"    text                               not null,
	enabled     boolean default true               not null,
	text        text                               not null,
	"channelId" text                               not null
		constraint "FK_15402461d2eb224e2e088b35606" references channels on update cascade on delete restrict,
	processed   boolean default false              not null,
	"isReply"   boolean default true               not null
);
create table if not exists channels_integrations
(
	id              uuid    default uuid_generate_v4() not null
		constraint "PK_d5f8b2ee43dbd8f5a68d1ce71ca" primary key,
	enabled         boolean default false              not null,
	"accessToken"   text,
	"refreshToken"  text,
	"clientId"      text,
	"clientSecret"  text,
	"apiKey"        text,
	data            jsonb,
	"channelId"     text                               not null
		constraint "FK_c17a48a983ac20ff800f553630a" references channels on update cascade on delete restrict,
	"integrationId" text                               not null
		constraint "FK_4958a98d1c19453c5755a422906" references integrations on update cascade on delete cascade
);
create table if not exists channels_keywords
(
	id                 uuid    default uuid_generate_v4() not null
		constraint "PK_6cb4cb2b8a95cc2728d1feab106" primary key,
	text               text                               not null,
	response           text,
	enabled            boolean default true               not null,
	cooldown           integer default 0                  not null,
	"channelId"        text                               not null
		constraint "FK_d17ea8fc949494c6b9681d5350c" references channels on update cascade on delete restrict,
	"cooldownExpireAt" timestamp,
	"isReply"          boolean default false              not null,
	"isRegular"        boolean default false              not null,
	usages             integer default 0                  not null
);
create unique index if not exists "channels_keywords_channelId_text_key" on channels_keywords ("channelId", text);
create table if not exists channels_moderation_settings
(
	id                   uuid    default uuid_generate_v4()     not null
		constraint "PK_fc3b021f2a7a5fde68b58543d93" primary key,
	type                 channels_moderation_settings_type_enum not null,
	enabled              boolean default false                  not null,
	subscribers          boolean default false                  not null,
	vips                 boolean default false                  not null,
	"banTime"            integer default 600                    not null,
	"banMessage"         text    default '' :: text             not null,
	"warningMessage"     text    default '' :: text             not null,
	"checkClips"         boolean default false,
	"triggerLength"      integer default 300,
	"maxPercentage"      integer default 50,
	"channelId"          text                                   not null
		constraint "FK_b34073350d4aa307b6104380e9b" references channels on update cascade on delete restrict,
	"blackListSentences" text[]  default '{}' :: text[]
);
create unique index if not exists "channels_moderation_settings_channelId_type_key" on channels_moderation_settings ("channelId", type);
create table if not exists channels_permits
(
	id          uuid default uuid_generate_v4() not null
		constraint "PK_1e24b77c350a387af97801f64ae" primary key,
	"channelId" text                            not null
		constraint "FK_b7e83e12e0482fa968de97b6b06" references channels on update cascade on delete restrict,
	"userId"    text                            not null
		constraint "FK_6bb136710060aa1be1744bc3bc9" references users on update cascade on delete restrict
);
create table if not exists channels_timers
(
	id                         uuid    default uuid_generate_v4() not null
		constraint "PK_858789b56e6e0956593768f9e05" primary key,
	name                       varchar(255)                       not null,
	enabled                    boolean default false              not null,
	"timeInterval"             integer default 0                  not null,
	"messageInterval"          integer default 0                  not null,
	"lastTriggerMessageNumber" integer default 0                  not null,
	"channelId"                text                               not null
		constraint "FK_50978b1848b99458d6ceb6b1989" references channels on update cascade on delete restrict
);
create table if not exists channels_dashboard_access
(
	id          text default gen_random_uuid() not null
		constraint "PK_33f482fb4aec2bed59de1d52cec" primary key,
	"channelId" text                           not null
		constraint "FK_82931b7f57c891fd7da375f3892" references channels on update cascade on delete restrict,
	"userId"    text                           not null
		constraint "FK_0ad4bccf63566cd4792dfe9f191" references users on update cascade on delete restrict
);
create table if not exists notifications
(
	id          text      default gen_random_uuid() not null
		constraint "PK_6a72c3c0f683f6462415e653c3a" primary key,
	"imageSrc"  text,
	"createdAt" timestamp default now()             not null,
	"userId"    text
		constraint "FK_692a909ee0fa9383e7859f9b406" references users on update cascade on delete
			set
			null
);
create table if not exists notifications_messages
(
	id               text default gen_random_uuid()       not null
		constraint "PK_10975a05eccf9793aa12dc143fa" primary key,
	text             text                                 not null,
	title            text,
	"langCode"       notifications_messages_langcode_enum not null,
	"notificationId" text                                 not null
		constraint "FK_4bfdb54ebcde220bcbcf696d862" references notifications on update cascade on delete restrict
);
create index if not exists "IDX_d98a275f8bc6cd986fcbe2eab0" on users ("tokenId");
create table if not exists users_files
(
	id       text default gen_random_uuid() not null
		constraint "PK_b01e2cd05bbd1ec1fddbd396ce3" primary key,
	name     text                           not null,
	size     integer                        not null,
	type     text                           not null,
	"userId" text                           not null
		constraint "FK_74cae0ea1fbb3df84c488ec0383" references users on update cascade on delete
			set
			null
);
create table if not exists users_stats
(
	id                  text    default gen_random_uuid() not null
		constraint "PK_44924448d5896c2364a4c6ddf75" primary key,
	messages            integer default 0                 not null,
	watched             bigint  default '0' :: bigint     not null,
	"channelId"         text                              not null
		constraint "FK_d55aab4a64c0c6b4b374b1da258" references channels on update cascade on delete restrict,
	"userId"            text                              not null
		constraint "FK_3d6cc217af2451426c44a30e678" references users on update cascade on delete restrict,
	"usedChannelPoints" bigint  default '0' :: bigint     not null
);
create unique index if not exists "users_stats_userId_channelId_key" on users_stats ("channelId", "userId");
create table if not exists users_viewed_notifications
(
	id               text      default gen_random_uuid() not null
		constraint "PK_8d45e86409cb2d99203c7b704bf" primary key,
	"createdAt"      timestamp default now()             not null,
	"notificationId" text                                not null
		constraint "FK_f5d19d90314d14d636752e2888b" references notifications on update cascade on delete restrict,
	"userId"         text                                not null
		constraint "FK_8d7e1a04a0d2d9868192561952d" references users on update cascade on delete restrict
);
create table if not exists channel_events_list
(
	id          uuid      default uuid_generate_v4() not null
		constraint "PK_ecc7dc4d42ccab404c461e25a4d" primary key,
	type        channel_events_list_type_enum        not null,
	"channelId" text                                 not null
		constraint "FK_62386183f7a7575141594880b03" references channels,
	"createdAt" timestamp default now()              not null
);
create table if not exists channel_events_follows
(
	id           uuid      default uuid_generate_v4() not null
		constraint "PK_0712d150127df2f1bc9d3ebd747" primary key,
	"eventId"    uuid                                 not null
		constraint "REL_4420e6ba1604749ff3184fefe0" unique
		constraint "FK_4420e6ba1604749ff3184fefe01" references channel_events_list,
	"fromUserId" varchar                              not null,
	"toUserId"   varchar                              not null,
	"createdAt"  timestamp default now()              not null
);
create table if not exists channel_events_donations
(
	id           uuid      default uuid_generate_v4() not null
		constraint "PK_968741787daf46977ee9e0f23c8" primary key,
	"eventId"    uuid                                 not null
		constraint "REL_e682e4bbfe3a354505755b7ee9" unique
		constraint "FK_e682e4bbfe3a354505755b7ee9b" references channel_events_list,
	"fromUserId" varchar,
	"toUserId"   varchar,
	amount       numeric                              not null,
	currency     varchar                              not null,
	username     varchar,
	message      varchar,
	"createdAt"  timestamp default now()              not null,
	"donateId"   varchar
		constraint "UQ_63834abea00e05711b7be7e3ff6" unique
);
create table if not exists channels_streams
(
	id               varchar           not null
		constraint "PK_d2b19f073bd39c68e9a9c6cca32" primary key,
	"userId"         text              not null
		constraint "FK_d2b9d6113cdeb816207be291ffa" references channels,
	"userLogin"      text              not null,
	"userName"       text              not null,
	"gameId"         integer           not null,
	"gameName"       text              not null,
	"communityIds"   text[]  default '{}':: text[],
	type             text              not null,
	title            text              not null,
	"viewerCount"    integer           not null,
	"startedAt"      timestamp         not null,
	language         text              not null,
	"thumbnailUrl"   text              not null,
	"tagIds"         text[]  default '{}':: text[],
	"isMature"       boolean           not null,
	"parsedMessages" integer default 0 not null,
	tags             text[]  default '{}':: text[]
);
create table if not exists users_online
(
	id          text default gen_random_uuid() not null
		constraint "PK_efb961b9e7ba28a837498e14827" primary key,
	"channelId" text                           not null
		constraint "FK_e6ae29713ab794b6ad8ef4fe5b4" references channels on update cascade on delete restrict,
	"userId"    text
		constraint "FK_e40473bd90abb17377f9dedb12a" references users on update cascade on delete restrict,
	"userName"  text
);
create table if not exists channels_moderation_warnings
(
	id          text default gen_random_uuid() not null
		constraint "PK_311c956bfd3e98c159ffb78740d" primary key,
	"channelId" text                           not null,
	"userId"    text                           not null,
	reason      text                           not null
);
create table if not exists channels_timers_responses
(
	id           uuid    default uuid_generate_v4() not null
		constraint "PK_9ab6dbe62d496cde6d8e52f11f9" primary key,
	text         text                               not null,
	"isAnnounce" boolean default true               not null,
	"timerId"    uuid                               not null
		constraint "FK_464cac4218d2ba9ac84325d06c8" references channels_timers on delete cascade
);
create table if not exists channels_requested_songs
(
	id                     uuid      default gen_random_uuid() not null
		constraint "PK_5f1d2b4311bd53a7cd499f1b59f" primary key,
	"videoId"              varchar                             not null,
	title                  text                                not null,
	duration               integer                             not null,
	"createdAt"            timestamp default now()             not null,
	"orderedById"          text                                not null
		constraint "FK_bf976a08dee5c9ffa7e1773defe" references users,
	"channelId"            text                                not null
		constraint "FK_a757e3014566676c024e4ce16d1" references channels,
	"orderedByName"        varchar                             not null,
	"deletedAt"            timestamp,
	"queuePosition"        integer                             not null,
	"orderedByDisplayName" varchar,
	"songLink"             varchar
);
create table if not exists channels_modules_settings
(
	id          uuid default gen_random_uuid()      not null
		constraint "PK_d5df4cbeb326c06be0e04654e36" primary key,
	type        channels_modules_settings_type_enum not null,
	settings    jsonb                               not null,
	"channelId" text                                not null
		constraint "FK_c145b2745bd936041f37b5d5d49" references channels,
	"userId"    text
		constraint "FK_b5f1c883e497ba7a0eeae08e8b8" references users,
	constraint "UQ_e74a3ef66bba62b18e3448211f7" unique ("channelId", "userId")
);
create table if not exists channels_messages
(
	"messageId"    text                    not null
		constraint "PK_685840c6efdbd345cf265976753" primary key,
	"channelId"    text                    not null
		constraint "FK_05a589d7e48f170714dc73243bf" references channels,
	"userId"       text                    not null
		constraint "FK_dd5560724c1166b9b70954f44be" references users,
	"userName"     text                    not null,
	text           text                    not null,
	"canBeDeleted" boolean   default true  not null,
	"createdAt"    timestamp default now() not null
);
create index if not exists "IDX_05a589d7e48f170714dc73243b" on channels_messages ("channelId");
create index if not exists "IDX_dd5560724c1166b9b70954f44b" on channels_messages ("userId");
create table if not exists users_ignored
(
	id            text not null
		constraint "PK_ac130ac65d82c39fee96465fb8d" primary key,
	login         text,
	"displayName" text
);
create index if not exists "IDX_ff17ac70a0b164274148218375" on users_ignored (login);
create index if not exists "IDX_74a54fd2a9d27143a7ed6bb979" on users_ignored ("displayName");
create table if not exists channels_emotes_usages
(
	id          uuid      default uuid_generate_v4() not null
		constraint "PK_cbaa7cf66062bb2a4d7926826f8" primary key,
	"channelId" text                                 not null
		constraint "FK_309ade49a31238d00065fc7c32e" references channels on update cascade on delete restrict,
	"userId"    text                                 not null
		constraint "FK_a736d0930dd0d26bd5f52bb0cf0" references users on update cascade on delete restrict,
	"createdAt" timestamp default now()              not null,
	emote       varchar                              not null
);
create table if not exists channels_events
(
	id          uuid    default uuid_generate_v4() not null
		constraint "PK_9e9ce619150ad05221318513d4d" primary key,
	type        channels_events_type_enum          not null,
	description text,
	"rewardId"  uuid,
	"commandId" text,
	"channelId" text                               not null
		constraint "FK_763ec88e86ecbf8ca6ec3a9ec7b" references channels,
	enabled     boolean default true               not null,
	"keywordId" text,
	online_only boolean default false              not null
);
create table if not exists channels_events_operations
(
	id               uuid    default uuid_generate_v4()   not null
		constraint "PK_52cef67c53dd212f1e1b3621a61" primary key,
	type             channels_events_operations_type_enum not null,
	delay            integer default 0                    not null,
	"eventId"        uuid                                 not null
		constraint "FK_b2e27e84fa5bfbf8fd27ac9e948" references channels_events,
	input            text,
	repeat           integer default 1                    not null,
	"order"          integer                              not null,
	"useAnnounce"    boolean default false                not null,
	"timeoutTime"    integer default 600                  not null,
	target           varchar,
	enabled          boolean default true                 not null,
	"timeoutMessage" text
);
create table if not exists channels_info_history
(
	id          uuid      default uuid_generate_v4() not null
		constraint "PK_deac5bf159894268f422f66c168" primary key,
	"channelId" text                                 not null
		constraint "FK_d326d1f33afefc8f45d6d546917" references channels,
	"createdAt" timestamp default now()              not null,
	title       text                                 not null,
	category    varchar                              not null
);
create table if not exists channels_commands_groups
(
	id          uuid    default uuid_generate_v4()                         not null
		constraint "PK_0e3eb1b93ec7980d98a87d4749a" primary key,
	"channelId" text                                                       not null
		constraint "FK_c202e2ed66394a1bd6651734078" references channels,
	name        varchar                                                    not null,
	color       varchar default 'rgba(37, 38, 43, 1)' :: character varying not null
);
create table if not exists channels_commands
(
	id                          uuid                                default uuid_generate_v4()                              not null
		constraint "PK_6f530d7cee75abfb1ef1cc7e09b" primary key,
	name                        text                                                                                        not null,
	cooldown                    integer                             default 0,
	"cooldownType"              channels_commands_cooldowntype_enum default 'GLOBAL' :: channels_commands_cooldowntype_enum not null,
	enabled                     boolean                             default true                                            not null,
	description                 text,
	visible                     boolean                             default true                                            not null,
	"default"                   boolean                             default false                                           not null,
	"defaultName"               text,
	module                      channels_commands_module_enum       default 'CUSTOM' :: channels_commands_module_enum       not null,
	"channelId"                 text                                                                                        not null
		constraint "FK_2ba87fd85e8e748470257452227" references channels on update cascade on delete restrict,
	aliases                     text[]                              default '{}' :: text[]                                  not null,
	is_reply                    boolean                             default true                                            not null,
	"keepResponsesOrder"        boolean                             default true                                            not null,
	"groupId"                   uuid
		constraint "FK_03d418a239ea80b80ebc562999b" references channels_commands_groups on delete
			set
			null,
	"rolesIds"                  text[]                              default '{}' :: text[]                                  not null,
	"deniedUsersIds"            text[]                              default '{}' :: text[]                                  not null,
	"allowedUsersIds"           text[]                              default '{}' :: text[]                                  not null,
	online_only                 boolean                             default false                                           not null,
	"requiredWatchTime"         integer                             default 0                                               not null,
	"requiredMessages"          integer                             default 0                                               not null,
	"requiredUsedChannelPoints" integer                             default 0                                               not null
);
create index if not exists "IDX_c340949fe9c2e1b1c636ff5ada" on channels_commands (name);
create index if not exists "IDX_2ba87fd85e8e74847025745222" on channels_commands ("channelId");
create unique index if not exists "channels_commands_name_channelId_key" on channels_commands ("channelId", name);
create table if not exists channels_commands_responses
(
	id          text    default gen_random_uuid() not null
		constraint "PK_2261731bc8ffc102ccc0082377b" primary key,
	text        text,
	"commandId" uuid                              not null
		constraint "FK_b18b2e298e41de0397acea55b97" references channels_commands on update cascade on delete cascade,
	"order"     integer default 0                 not null
);
create table if not exists channels_commands_usages
(
	id          text default gen_random_uuid() not null
		constraint "PK_e7549bc5ec5f2f2d779102caa8a" primary key,
	"channelId" text                           not null,
	"commandId" uuid                           not null
		constraint "FK_2db68a1b263bea3f214186bb24c" references channels_commands on update cascade on delete cascade,
	"userId"    text                           not null
		constraint "FK_14be9e16eaece94af10d56044e0" references users on update cascade on delete restrict
);
create table if not exists channels_roles
(
	id          uuid                              default uuid_generate_v4()                        not null
		constraint "PK_8445a7f7c362bbebcaab4de9845" primary key,
	"channelId" text                                                                                not null
		constraint "FK_1c6f5f58e54b77d7480a4895103" references channels,
	name        varchar                                                                             not null,
	type        channels_roles_type_enum          default 'CUSTOM' :: channels_roles_type_enum      not null,
	permissions channels_roles_permissions_enum[] default '{}' :: channels_roles_permissions_enum[] not null,
	settings    jsonb                             default '{}' :: jsonb                             not null
);
create table if not exists channels_roles_users
(
	id       uuid default uuid_generate_v4() not null
		constraint "PK_22c498bd2efb082200bb30ad44e" primary key,
	"userId" text                            not null
		constraint "FK_c3b3d3917f24ec3fe3f5049e667" references users,
	"roleId" uuid                            not null
		constraint "FK_ae2e54203d7fed7ed38243e8b13" references channels_roles on delete cascade
);
create table if not exists channels_events_operations_filters
(
	id            uuid default uuid_generate_v4()              not null
		constraint "PK_184c3468090c81affe77014f162" primary key,
	"operationId" uuid                                         not null
		constraint "FK_d56a4ca65d44fc4adbaf66a2e80" references channels_events_operations on delete cascade,
	type          channels_events_operations_filters_type_enum not null,
	"left"        text                                         not null,
	"right"       text                                         not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
