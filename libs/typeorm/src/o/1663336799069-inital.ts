import { MigrationInterface, QueryRunner } from 'typeorm';

export class inital1663336799069 implements MigrationInterface {
  name = 'inital1663336799069';

  public async up(queryRunner: QueryRunner): Promise<void> {
    // check if table exists
    const table = await queryRunner.getTable('bots');
    if (table) {
      return;
    }

    await queryRunner.query(`
    CREATE TABLE "public"."_prisma_migrations" (
        "id" character varying(36) NOT NULL,
        "checksum" character varying(64) NOT NULL,
        "finished_at" timestamptz,
        "migration_name" character varying(255) NOT NULL,
        "logs" text,
        "rolled_back_at" timestamptz,
        "started_at" timestamptz DEFAULT now() NOT NULL,
        "applied_steps_count" integer DEFAULT '0' NOT NULL,
        CONSTRAINT "_prisma_migrations_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false); 
    
    CREATE TABLE "public"."bots" (
        "id" text NOT NULL,
        "type" "BotType" NOT NULL,
        "tokenId" text,
        CONSTRAINT "bots_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "bots_tokenId_key" UNIQUE ("tokenId")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels" (
        "id" text NOT NULL,
        "isEnabled" boolean DEFAULT true NOT NULL,
        "isTwitchBanned" boolean DEFAULT false NOT NULL,
        "isBanned" boolean DEFAULT false NOT NULL,
        "botId" text NOT NULL,
        CONSTRAINT "channels_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_commands" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "name" text NOT NULL,
        "cooldown" integer DEFAULT '0',
        "cooldownType" "CooldownType" DEFAULT GLOBAL NOT NULL,
        "enabled" boolean DEFAULT true NOT NULL,
        "aliases" jsonb DEFAULT '[]',
        "description" text,
        "visible" boolean DEFAULT true NOT NULL,
        "channelId" text NOT NULL,
        "permission" "CommandPermission" NOT NULL,
        "default" boolean DEFAULT false NOT NULL,
        "defaultName" text,
        "module" "CommandModule" DEFAULT CUSTOM NOT NULL,
        CONSTRAINT "channels_commands_name_channelId_key" UNIQUE ("name", "channelId"),
        CONSTRAINT "channels_commands_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    CREATE INDEX "channels_commands_channelId_idx" ON "public"."channels_commands" USING btree ("channelId");
    
    CREATE INDEX "channels_commands_name_idx" ON "public"."channels_commands" USING btree ("name");
    
    
    CREATE TABLE "public"."channels_commands_responses" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "text" text,
        "commandId" text NOT NULL,
        CONSTRAINT "channels_commands_responses_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_commands_usages" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "userId" text NOT NULL,
        "channelId" text NOT NULL,
        "commandId" text NOT NULL,
        CONSTRAINT "channels_commands_usages_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_customvars" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "name" text NOT NULL,
        "description" text,
        "type" "CustomVarType" NOT NULL,
        "evalValue" text,
        "response" text,
        "channelId" text NOT NULL,
        CONSTRAINT "channels_customvars_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_dashboard_access" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "channelId" text NOT NULL,
        "userId" text NOT NULL,
        CONSTRAINT "channels_dashboard_access_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_dota_accounts" (
        "id" text NOT NULL,
        "channelId" text NOT NULL,
        CONSTRAINT "channels_dota_accounts_id_channelId_key" UNIQUE ("id", "channelId"),
        CONSTRAINT "channels_dota_accounts_pkey" PRIMARY KEY ("id", "channelId")
    ) WITH (oids = false);
    
    CREATE INDEX "channels_dota_accounts_id_idx" ON "public"."channels_dota_accounts" USING btree ("id");
    
    
    CREATE TABLE "public"."channels_greetings" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "channelId" text NOT NULL,
        "userId" text NOT NULL,
        "enabled" boolean DEFAULT true NOT NULL,
        "text" text NOT NULL,
        CONSTRAINT "channels_greetings_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_integrations" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "enabled" boolean DEFAULT false NOT NULL,
        "channelId" text NOT NULL,
        "integrationId" text NOT NULL,
        "accessToken" text,
        "refreshToken" text,
        "clientId" text,
        "clientSecret" text,
        "apiKey" text,
        "data" jsonb,
        CONSTRAINT "channels_integrations_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_keywords" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "channelId" text NOT NULL,
        "text" text NOT NULL,
        "response" text NOT NULL,
        "enabled" boolean DEFAULT true NOT NULL,
        "cooldown" integer DEFAULT '0',
        CONSTRAINT "channels_keywords_channelId_text_key" UNIQUE ("channelId", "text"),
        CONSTRAINT "channels_keywords_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_moderation_settings" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "type" "SettingsType" NOT NULL,
        "channelId" text NOT NULL,
        "enabled" boolean DEFAULT false NOT NULL,
        "subscribers" boolean DEFAULT false NOT NULL,
        "vips" boolean DEFAULT false NOT NULL,
        "banTime" integer DEFAULT '600' NOT NULL,
        "banMessage" text,
        "warningMessage" text,
        "checkClips" boolean DEFAULT false,
        "triggerLength" integer DEFAULT '300',
        "maxPercentage" integer DEFAULT '50',
        "blackListSentences" jsonb DEFAULT '[]',
        CONSTRAINT "channels_moderation_settings_channelId_type_key" UNIQUE ("channelId", "type"),
        CONSTRAINT "channels_moderation_settings_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_permits" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "channelId" text NOT NULL,
        "userId" text NOT NULL,
        CONSTRAINT "channels_permits_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."channels_timers" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "channelId" text NOT NULL,
        "name" character varying(255) NOT NULL,
        "enabled" boolean DEFAULT true NOT NULL,
        "responses" jsonb DEFAULT '[]' NOT NULL,
        "last" integer DEFAULT '0' NOT NULL,
        "timeInterval" integer DEFAULT '0' NOT NULL,
        "messageInterval" integer DEFAULT '0' NOT NULL,
        "lastTriggerMessageNumber" integer DEFAULT '0' NOT NULL,
        CONSTRAINT "channels_timers_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."dota_game_modes" (
        "id" integer NOT NULL,
        "name" text NOT NULL,
        CONSTRAINT "dota_game_modes_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."dota_heroes" (
        "id" integer NOT NULL,
        "name" text NOT NULL,
        CONSTRAINT "dota_heroes_id_key" UNIQUE ("id"),
        CONSTRAINT "dota_heroes_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."dota_matches" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "startedAt" timestamp(3) NOT NULL,
        "lobby_type" integer NOT NULL,
        "gameModeId" integer NOT NULL,
        "players" integer[],
        "players_heroes" integer[],
        "weekend_tourney_bracket_round" text,
        "weekend_tourney_skill_level" text,
        "match_id" text NOT NULL,
        "avarage_mmr" integer NOT NULL,
        "lobbyId" text NOT NULL,
        "finished" boolean DEFAULT false NOT NULL,
        CONSTRAINT "dota_matches_match_id_key" UNIQUE ("match_id"),
        CONSTRAINT "dota_matches_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."dota_matches_cards" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "match_id" text NOT NULL,
        "account_id" text NOT NULL,
        "rank_tier" integer,
        "leaderboard_rank" integer,
        CONSTRAINT "dota_matches_cards_match_id_account_id_key" UNIQUE ("match_id", "account_id"),
        CONSTRAINT "dota_matches_cards_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."dota_matches_results" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "match_id" text NOT NULL,
        "players" jsonb NOT NULL,
        "radiant_win" boolean NOT NULL,
        "game_mode" integer NOT NULL,
        CONSTRAINT "dota_matches_results_match_id_key" UNIQUE ("match_id"),
        CONSTRAINT "dota_matches_results_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."dota_medals" (
        "rank_tier" text NOT NULL,
        "name" text NOT NULL,
        CONSTRAINT "dota_medals_pkey" PRIMARY KEY ("rank_tier")
    ) WITH (oids = false);
    
    CREATE INDEX "dota_medals_rank_tier_idx" ON "public"."dota_medals" USING btree ("rank_tier");
    
    
    CREATE TABLE "public"."integrations" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "service" "IntegrationService" NOT NULL,
        "accessToken" text,
        "refreshToken" text,
        "clientId" text,
        "clientSecret" text,
        "apiKey" text,
        "redirectUrl" text,
        CONSTRAINT "integrations_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."notifications" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "imageSrc" text,
        "createdAt" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
        "userId" text,
        CONSTRAINT "notifications_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."notifications_messages" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "text" text NOT NULL,
        "title" text,
        "langCode" "LangCode" NOT NULL,
        "notificationId" text NOT NULL,
        CONSTRAINT "notifications_messages_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."tokens" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "accessToken" text NOT NULL,
        "refreshToken" text NOT NULL,
        "expiresIn" integer NOT NULL,
        "obtainmentTimestamp" timestamp(3) NOT NULL,
        CONSTRAINT "tokens_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."users" (
        "id" text NOT NULL,
        "tokenId" text,
        "isTester" boolean DEFAULT false NOT NULL,
        "isBotAdmin" boolean DEFAULT false NOT NULL,
        CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "users_tokenId_key" UNIQUE ("tokenId")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."users_files" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "name" text NOT NULL,
        "size" integer NOT NULL,
        "type" text NOT NULL,
        "userId" text,
        CONSTRAINT "users_files_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."users_stats" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "userId" text NOT NULL,
        "channelId" text NOT NULL,
        "messages" integer DEFAULT '0' NOT NULL,
        "watched" bigint DEFAULT '0' NOT NULL,
        CONSTRAINT "users_stats_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "users_stats_userId_channelId_key" UNIQUE ("userId", "channelId")
    ) WITH (oids = false);
    
    
    CREATE TABLE "public"."users_viewed_notifications" (
        "id" text DEFAULT 'gen_random_uuid()' NOT NULL,
        "userId" text NOT NULL,
        "notificationId" text NOT NULL,
        "createdAt" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
        CONSTRAINT "users_viewed_notifications_pkey" PRIMARY KEY ("id")
    ) WITH (oids = false);
    
    
    ALTER TABLE ONLY "public"."bots" ADD CONSTRAINT "bots_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES tokens(id) ON UPDATE CASCADE ON DELETE SET NULL NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels" ADD CONSTRAINT "channels_botId_fkey" FOREIGN KEY ("botId") REFERENCES bots(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."channels" ADD CONSTRAINT "channels_id_fkey" FOREIGN KEY (id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_commands" ADD CONSTRAINT "channels_commands_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_commands_responses" ADD CONSTRAINT "channels_commands_responses_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES channels_commands(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES channels_commands(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_customvars" ADD CONSTRAINT "channels_customvars_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_dashboard_access" ADD CONSTRAINT "channels_dashboard_access_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."channels_dashboard_access" ADD CONSTRAINT "channels_dashboard_access_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_dota_accounts" ADD CONSTRAINT "channels_dota_accounts_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_greetings" ADD CONSTRAINT "channels_greetings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_integrations" ADD CONSTRAINT "channels_integrations_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."channels_integrations" ADD CONSTRAINT "channels_integrations_integrationId_fkey" FOREIGN KEY ("integrationId") REFERENCES integrations(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_keywords" ADD CONSTRAINT "channels_keywords_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_moderation_settings" ADD CONSTRAINT "channels_moderation_settings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_permits" ADD CONSTRAINT "channels_permits_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."channels_permits" ADD CONSTRAINT "channels_permits_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."channels_timers" ADD CONSTRAINT "channels_timers_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."dota_matches" ADD CONSTRAINT "dota_matches_gameModeId_fkey" FOREIGN KEY ("gameModeId") REFERENCES dota_game_modes(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."dota_matches_cards" ADD CONSTRAINT "dota_matches_cards_match_id_fkey" FOREIGN KEY (match_id) REFERENCES dota_matches(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."dota_matches_results" ADD CONSTRAINT "dota_matches_results_match_id_fkey" FOREIGN KEY (match_id) REFERENCES dota_matches(match_id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."notifications" ADD CONSTRAINT "notifications_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."notifications_messages" ADD CONSTRAINT "notifications_messages_notificationId_fkey" FOREIGN KEY ("notificationId") REFERENCES notifications(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES tokens(id) ON UPDATE CASCADE ON DELETE SET NULL NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."users_files" ADD CONSTRAINT "users_files_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."users_stats" ADD CONSTRAINT "users_stats_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."users_stats" ADD CONSTRAINT "users_stats_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    
    ALTER TABLE ONLY "public"."users_viewed_notifications" ADD CONSTRAINT "users_viewed_notifications_notificationId_fkey" FOREIGN KEY ("notificationId") REFERENCES notifications(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    ALTER TABLE ONLY "public"."users_viewed_notifications" ADD CONSTRAINT "users_viewed_notifications_userId_fkey" FOREIGN KEY ("userId") REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
    `);
  }
  down(queryRunner: QueryRunner): Promise<any> {
    throw new Error('Method not implemented.');
  }
}
