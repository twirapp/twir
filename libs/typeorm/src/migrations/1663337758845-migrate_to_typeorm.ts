import { MigrationInterface, QueryRunner } from 'typeorm';

export class migrateToTypeorm1663337758845 implements MigrationInterface {
  name = 'migrateToTypeorm1663337758845';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "bots" DROP CONSTRAINT "bots_tokenId_fkey"`);
    await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "channels_botId_fkey"`);
    await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "channels_id_fkey"`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" DROP CONSTRAINT "channels_commands_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" DROP CONSTRAINT "channels_customvars_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "channels_dota_accounts_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" DROP CONSTRAINT "channels_greetings_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" DROP CONSTRAINT "channels_integrations_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" DROP CONSTRAINT "channels_integrations_integrationId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" DROP CONSTRAINT "channels_keywords_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" DROP CONSTRAINT "channels_moderation_settings_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" DROP CONSTRAINT "channels_permits_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" DROP CONSTRAINT "channels_permits_userId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" DROP CONSTRAINT "channels_timers_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" DROP CONSTRAINT "channels_commands_responses_commandId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" DROP CONSTRAINT "channels_commands_usages_userId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" DROP CONSTRAINT "channels_commands_usages_commandId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" DROP CONSTRAINT "channels_dashboard_access_channelId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" DROP CONSTRAINT "channels_dashboard_access_userId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" DROP CONSTRAINT "dota_matches_gameModeId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_cards" DROP CONSTRAINT "dota_matches_cards_match_id_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" DROP CONSTRAINT "dota_matches_results_match_id_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" DROP CONSTRAINT "notifications_userId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" DROP CONSTRAINT "notifications_messages_notificationId_fkey"`,
    );
    await queryRunner.query(`ALTER TABLE "users" DROP CONSTRAINT "users_tokenId_fkey"`);
    await queryRunner.query(`ALTER TABLE "users_files" DROP CONSTRAINT "users_files_userId_fkey"`);
    await queryRunner.query(
      `ALTER TABLE "users_stats" DROP CONSTRAINT "users_stats_channelId_fkey"`,
    );
    await queryRunner.query(`ALTER TABLE "users_stats" DROP CONSTRAINT "users_stats_userId_fkey"`);
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" DROP CONSTRAINT "users_viewed_notifications_notificationId_fkey"`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" DROP CONSTRAINT "users_viewed_notifications_userId_fkey"`,
    );
    await queryRunner.query(`DROP INDEX "public"."bots_tokenId_key"`);
    await queryRunner.query(`DROP INDEX "public"."channels_commands_channelId_idx"`);
    await queryRunner.query(`DROP INDEX "public"."channels_commands_name_idx"`);
    await queryRunner.query(`DROP INDEX "public"."channels_dota_accounts_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."dota_heroes_id_key"`);
    await queryRunner.query(`DROP INDEX "public"."dota_matches_match_id_key"`);
    await queryRunner.query(`DROP INDEX "public"."dota_medals_rank_tier_idx"`);
    await queryRunner.query(`DROP INDEX "public"."users_tokenId_key"`);
    await queryRunner.query(`ALTER TYPE "public"."BotType" RENAME TO "BotType_old"`);
    await queryRunner.query(`CREATE TYPE "public"."bots_type_enum" AS ENUM('DEFAULT', 'CUSTOM')`);
    await queryRunner.query(
      `ALTER TABLE "bots" ALTER COLUMN "type" TYPE "public"."bots_type_enum" USING "type"::"text"::"public"."bots_type_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."BotType_old"`);
    await queryRunner.query(
      `ALTER TABLE "bots" ADD CONSTRAINT "UQ_df4240a5d71aa6a23b829d3cee8" UNIQUE ("tokenId")`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" ADD CONSTRAINT "UQ_bc603823f3f741359c2339389f9" UNIQUE ("id")`,
    );
    await queryRunner.query(`ALTER TABLE "channels" ALTER COLUMN "botId" DROP NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(`ALTER TYPE "public"."CooldownType" RENAME TO "CooldownType_old"`);
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_cooldowntype_enum" AS ENUM('GLOBAL', 'PER_USER')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "cooldownType" DROP DEFAULT`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "cooldownType" TYPE "public"."channels_commands_cooldowntype_enum" USING "cooldownType"::"text"::"public"."channels_commands_cooldowntype_enum"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "cooldownType" SET DEFAULT 'GLOBAL'`,
    );
    await queryRunner.query(`DROP TYPE "public"."CooldownType_old"`);
    await queryRunner.query(
      `ALTER TYPE "public"."CommandPermission" RENAME TO "CommandPermission_old"`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_permission_enum" AS ENUM('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "permission" TYPE "public"."channels_commands_permission_enum" USING "permission"::"text"::"public"."channels_commands_permission_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."CommandPermission_old"`);
    await queryRunner.query(`ALTER TYPE "public"."CommandModule" RENAME TO "CommandModule_old"`);
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_module_enum" AS ENUM('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION')`,
    );
    await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum" USING "module"::"text"::"public"."channels_commands_module_enum"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`,
    );
    await queryRunner.query(`DROP TYPE "public"."CommandModule_old"`);
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(`ALTER TYPE "public"."CustomVarType" RENAME TO "CustomVarType_old"`);
    await queryRunner.query(
      `CREATE TYPE "public"."channels_customvars_type_enum" AS ENUM('SCRIPT', 'TEXT')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ALTER COLUMN "type" TYPE "public"."channels_customvars_type_enum" USING "type"::"text"::"public"."channels_customvars_type_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."CustomVarType_old"`);
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ALTER COLUMN "channelId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ALTER COLUMN "channelId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ALTER COLUMN "channelId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ALTER COLUMN "integrationId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."channels_moderation_settings_channelId_type_key"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(`ALTER TYPE "public"."SettingsType" RENAME TO "SettingsType_old"`);
    await queryRunner.query(
      `CREATE TYPE "public"."channels_moderation_settings_type_enum" AS ENUM('links', 'blacklists', 'symbols', 'longMessage', 'caps', 'emotes')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "type" TYPE "public"."channels_moderation_settings_type_enum" USING "type"::"text"::"public"."channels_moderation_settings_type_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."SettingsType_old"`);
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ALTER COLUMN "channelId" DROP NOT NULL`,
    );
    await queryRunner.query(`ALTER TABLE "channels_permits" ALTER COLUMN "userId" DROP NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ALTER COLUMN "enabled" SET DEFAULT false`,
    );
    await queryRunner.query(`ALTER TABLE "channels_timers" ALTER COLUMN "channelId" DROP NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ALTER COLUMN "commandId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ALTER COLUMN "commandId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ALTER COLUMN "userId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ALTER COLUMN "channelId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ALTER COLUMN "userId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(`ALTER TABLE "dota_matches" ALTER COLUMN "startedAt" TYPE TIMESTAMP`);
    await queryRunner.query(
      `ALTER TABLE "dota_matches" ADD CONSTRAINT "UQ_27d39ee5a1cdb04ed78a52286a1" UNIQUE ("match_id")`,
    );
    await queryRunner.query(`ALTER TABLE "dota_matches" ALTER COLUMN "gameModeId" DROP NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "dota_matches_cards" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" ADD CONSTRAINT "UQ_66717d38c93b88a283ad77736ad" UNIQUE ("match_id")`,
    );
    await queryRunner.query(
      `ALTER TABLE "integrations" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TYPE "public"."IntegrationService" RENAME TO "IntegrationService_old"`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS')`,
    );
    await queryRunner.query(
      `ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."IntegrationService_old"`);
    await queryRunner.query(
      `ALTER TABLE "notifications" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(`ALTER TABLE "notifications" ALTER COLUMN "createdAt" TYPE TIMESTAMP`);
    await queryRunner.query(
      `ALTER TABLE "notifications" ALTER COLUMN "createdAt" SET DEFAULT now()`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(`ALTER TYPE "public"."LangCode" RENAME TO "LangCode_old"`);
    await queryRunner.query(
      `CREATE TYPE "public"."notifications_messages_langcode_enum" AS ENUM('RU', 'GB')`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ALTER COLUMN "langCode" TYPE "public"."notifications_messages_langcode_enum" USING "langCode"::"text"::"public"."notifications_messages_langcode_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."LangCode_old"`);
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ALTER COLUMN "notificationId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "tokens" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "tokens" ALTER COLUMN "obtainmentTimestamp" TYPE TIMESTAMP`,
    );
    await queryRunner.query(
      `ALTER TABLE "users" ADD CONSTRAINT "UQ_d98a275f8bc6cd986fcbe2eab01" UNIQUE ("tokenId")`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_files" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "id" SET DEFAULT 'gen_random_uuid()'`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "createdAt" TYPE TIMESTAMP`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "createdAt" SET DEFAULT now()`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "notificationId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "userId" DROP NOT NULL`,
    );
    await queryRunner.query(`CREATE INDEX "IDX_df4240a5d71aa6a23b829d3cee" ON "bots" ("tokenId") `);
    await queryRunner.query(
      `CREATE INDEX "IDX_c340949fe9c2e1b1c636ff5ada" ON "channels_commands" ("name") `,
    );
    await queryRunner.query(
      `CREATE INDEX "IDX_2ba87fd85e8e74847025745222" ON "channels_commands" ("channelId") `,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "channels_moderation_settings_channelId_type_key" ON "channels_moderation_settings" ("channelId", "type") `,
    );
    await queryRunner.query(
      `CREATE INDEX "IDX_82490264788f0b786e2ff312a5" ON "dota_heroes" ("name") `,
    );
    await queryRunner.query(
      `CREATE INDEX "IDX_27d39ee5a1cdb04ed78a52286a" ON "dota_matches" ("match_id") `,
    );
    await queryRunner.query(
      `CREATE INDEX "IDX_e32891d3f8e343fe1d0107fa4e" ON "dota_medals" ("name") `,
    );
    await queryRunner.query(
      `CREATE INDEX "IDX_d98a275f8bc6cd986fcbe2eab0" ON "users" ("tokenId") `,
    );
    await queryRunner.query(
      `ALTER TABLE "bots" ADD CONSTRAINT "bots_tokenId_key" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" ADD CONSTRAINT "FK_4f890144c0cb55fe7867b8f61e6" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" ADD CONSTRAINT "FK_bc603823f3f741359c2339389f9" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_2ba87fd85e8e748470257452227" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ADD CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ADD CONSTRAINT "FK_15402461d2eb224e2e088b35606" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ADD CONSTRAINT "FK_c17a48a983ac20ff800f553630a" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ADD CONSTRAINT "FK_4958a98d1c19453c5755a422906" FOREIGN KEY ("integrationId") REFERENCES "integrations"("id") ON DELETE CASCADE ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ADD CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "FK_b34073350d4aa307b6104380e9b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ADD CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ADD CONSTRAINT "FK_6bb136710060aa1be1744bc3bc9" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ADD CONSTRAINT "FK_50978b1848b99458d6ceb6b1989" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ADD CONSTRAINT "FK_b18b2e298e41de0397acea55b97" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE CASCADE ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "FK_2db68a1b263bea3f214186bb24c" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE CASCADE ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "FK_14be9e16eaece94af10d56044e0" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "FK_82931b7f57c891fd7da375f3892" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "FK_0ad4bccf63566cd4792dfe9f191" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" ADD CONSTRAINT "FK_7fa32afcc462ae6f50b176f6c9c" FOREIGN KEY ("gameModeId") REFERENCES "dota_game_modes"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_cards" ADD CONSTRAINT "FK_a9d58830e97acdda53124a0431f" FOREIGN KEY ("match_id") REFERENCES "dota_matches"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" ADD CONSTRAINT "FK_66717d38c93b88a283ad77736ad" FOREIGN KEY ("match_id") REFERENCES "dota_matches"("match_id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" ADD CONSTRAINT "FK_692a909ee0fa9383e7859f9b406" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ADD CONSTRAINT "FK_4bfdb54ebcde220bcbcf696d862" FOREIGN KEY ("notificationId") REFERENCES "notifications"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users" ADD CONSTRAINT "FK_d98a275f8bc6cd986fcbe2eab01" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_files" ADD CONSTRAINT "FK_74cae0ea1fbb3df84c488ec0383" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" ADD CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" ADD CONSTRAINT "FK_3d6cc217af2451426c44a30e678" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ADD CONSTRAINT "FK_f5d19d90314d14d636752e2888b" FOREIGN KEY ("notificationId") REFERENCES "notifications"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ADD CONSTRAINT "FK_8d7e1a04a0d2d9868192561952d" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" DROP CONSTRAINT "FK_8d7e1a04a0d2d9868192561952d"`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" DROP CONSTRAINT "FK_f5d19d90314d14d636752e2888b"`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" DROP CONSTRAINT "FK_3d6cc217af2451426c44a30e678"`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" DROP CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258"`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_files" DROP CONSTRAINT "FK_74cae0ea1fbb3df84c488ec0383"`,
    );
    await queryRunner.query(`ALTER TABLE "users" DROP CONSTRAINT "FK_d98a275f8bc6cd986fcbe2eab01"`);
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" DROP CONSTRAINT "FK_4bfdb54ebcde220bcbcf696d862"`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" DROP CONSTRAINT "FK_692a909ee0fa9383e7859f9b406"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" DROP CONSTRAINT "FK_66717d38c93b88a283ad77736ad"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_cards" DROP CONSTRAINT "FK_a9d58830e97acdda53124a0431f"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" DROP CONSTRAINT "FK_7fa32afcc462ae6f50b176f6c9c"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" DROP CONSTRAINT "FK_0ad4bccf63566cd4792dfe9f191"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" DROP CONSTRAINT "FK_82931b7f57c891fd7da375f3892"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" DROP CONSTRAINT "FK_14be9e16eaece94af10d56044e0"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" DROP CONSTRAINT "FK_2db68a1b263bea3f214186bb24c"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" DROP CONSTRAINT "FK_b18b2e298e41de0397acea55b97"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" DROP CONSTRAINT "FK_50978b1848b99458d6ceb6b1989"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" DROP CONSTRAINT "FK_6bb136710060aa1be1744bc3bc9"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" DROP CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" DROP CONSTRAINT "FK_b34073350d4aa307b6104380e9b"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" DROP CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" DROP CONSTRAINT "FK_4958a98d1c19453c5755a422906"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" DROP CONSTRAINT "FK_c17a48a983ac20ff800f553630a"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" DROP CONSTRAINT "FK_15402461d2eb224e2e088b35606"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" DROP CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_2ba87fd85e8e748470257452227"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" DROP CONSTRAINT "FK_bc603823f3f741359c2339389f9"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" DROP CONSTRAINT "FK_4f890144c0cb55fe7867b8f61e6"`,
    );
    await queryRunner.query(`ALTER TABLE "bots" DROP CONSTRAINT "bots_tokenId_key"`);
    await queryRunner.query(`DROP INDEX "public"."IDX_d98a275f8bc6cd986fcbe2eab0"`);
    await queryRunner.query(`DROP INDEX "public"."IDX_e32891d3f8e343fe1d0107fa4e"`);
    await queryRunner.query(`DROP INDEX "public"."IDX_27d39ee5a1cdb04ed78a52286a"`);
    await queryRunner.query(`DROP INDEX "public"."IDX_82490264788f0b786e2ff312a5"`);
    await queryRunner.query(
      `DROP INDEX "public"."channels_moderation_settings_channelId_type_key"`,
    );
    await queryRunner.query(`DROP INDEX "public"."IDX_2ba87fd85e8e74847025745222"`);
    await queryRunner.query(`DROP INDEX "public"."IDX_c340949fe9c2e1b1c636ff5ada"`);
    await queryRunner.query(`DROP INDEX "public"."IDX_df4240a5d71aa6a23b829d3cee"`);
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "userId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "notificationId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "createdAt" SET DEFAULT CURRENT_TIMESTAMP`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "createdAt" TYPE TIMESTAMP(3)`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_files" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(`ALTER TABLE "users" DROP CONSTRAINT "UQ_d98a275f8bc6cd986fcbe2eab01"`);
    await queryRunner.query(
      `ALTER TABLE "tokens" ALTER COLUMN "obtainmentTimestamp" TYPE TIMESTAMP(3)`,
    );
    await queryRunner.query(`ALTER TABLE "tokens" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ALTER COLUMN "notificationId" SET NOT NULL`,
    );
    await queryRunner.query(`CREATE TYPE "public"."LangCode_old" AS ENUM('RU', 'GB')`);
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ALTER COLUMN "langCode" TYPE "public"."LangCode_old" USING "langCode"::"text"::"public"."LangCode_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."notifications_messages_langcode_enum"`);
    await queryRunner.query(`ALTER TYPE "public"."LangCode_old" RENAME TO "LangCode"`);
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" ALTER COLUMN "createdAt" SET DEFAULT CURRENT_TIMESTAMP`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" ALTER COLUMN "createdAt" TYPE TIMESTAMP(3)`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."IntegrationService_old" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS')`,
    );
    await queryRunner.query(
      `ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."IntegrationService_old" USING "service"::"text"::"public"."IntegrationService_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
    await queryRunner.query(
      `ALTER TYPE "public"."IntegrationService_old" RENAME TO "IntegrationService"`,
    );
    await queryRunner.query(
      `ALTER TABLE "integrations" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" DROP CONSTRAINT "UQ_66717d38c93b88a283ad77736ad"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_cards" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(`ALTER TABLE "dota_matches" ALTER COLUMN "gameModeId" SET NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "dota_matches" DROP CONSTRAINT "UQ_27d39ee5a1cdb04ed78a52286a1"`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" ALTER COLUMN "startedAt" TYPE TIMESTAMP(3)`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ALTER COLUMN "userId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ALTER COLUMN "channelId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ALTER COLUMN "userId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ALTER COLUMN "commandId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ALTER COLUMN "commandId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(`ALTER TABLE "channels_timers" ALTER COLUMN "channelId" SET NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ALTER COLUMN "enabled" SET DEFAULT true`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(`ALTER TABLE "channels_permits" ALTER COLUMN "userId" SET NOT NULL`);
    await queryRunner.query(`ALTER TABLE "channels_permits" ALTER COLUMN "channelId" SET NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."SettingsType_old" AS ENUM('links', 'blacklists', 'symbols', 'longMessage', 'caps', 'emotes')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "type" TYPE "public"."SettingsType_old" USING "type"::"text"::"public"."SettingsType_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_moderation_settings_type_enum"`);
    await queryRunner.query(`ALTER TYPE "public"."SettingsType_old" RENAME TO "SettingsType"`);
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "channels_moderation_settings_channelId_type_key" ON "channels_moderation_settings" ("type", "channelId") `,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ALTER COLUMN "integrationId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ALTER COLUMN "channelId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ALTER COLUMN "channelId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ALTER COLUMN "channelId" SET NOT NULL`,
    );
    await queryRunner.query(`CREATE TYPE "public"."CustomVarType_old" AS ENUM('SCRIPT', 'TEXT')`);
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ALTER COLUMN "type" TYPE "public"."CustomVarType_old" USING "type"::"text"::"public"."CustomVarType_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_customvars_type_enum"`);
    await queryRunner.query(`ALTER TYPE "public"."CustomVarType_old" RENAME TO "CustomVarType"`);
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."CommandModule_old" AS ENUM('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION')`,
    );
    await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."CommandModule_old" USING "module"::"text"::"public"."CommandModule_old"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum"`);
    await queryRunner.query(`ALTER TYPE "public"."CommandModule_old" RENAME TO "CommandModule"`);
    await queryRunner.query(
      `CREATE TYPE "public"."CommandPermission_old" AS ENUM('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "permission" TYPE "public"."CommandPermission_old" USING "permission"::"text"::"public"."CommandPermission_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_permission_enum"`);
    await queryRunner.query(
      `ALTER TYPE "public"."CommandPermission_old" RENAME TO "CommandPermission"`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."CooldownType_old" AS ENUM('GLOBAL', 'PER_USER')`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "cooldownType" DROP DEFAULT`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "cooldownType" TYPE "public"."CooldownType_old" USING "cooldownType"::"text"::"public"."CooldownType_old"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "cooldownType" SET DEFAULT 'GLOBAL'`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_cooldowntype_enum"`);
    await queryRunner.query(`ALTER TYPE "public"."CooldownType_old" RENAME TO "CooldownType"`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(`ALTER TABLE "channels" ALTER COLUMN "botId" SET NOT NULL`);
    await queryRunner.query(
      `ALTER TABLE "channels" DROP CONSTRAINT "UQ_bc603823f3f741359c2339389f9"`,
    );
    await queryRunner.query(`ALTER TABLE "bots" DROP CONSTRAINT "UQ_df4240a5d71aa6a23b829d3cee8"`);
    await queryRunner.query(`CREATE TYPE "public"."BotType_old" AS ENUM('DEFAULT', 'CUSTOM')`);
    await queryRunner.query(
      `ALTER TABLE "bots" ALTER COLUMN "type" TYPE "public"."BotType_old" USING "type"::"text"::"public"."BotType_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."bots_type_enum"`);
    await queryRunner.query(`ALTER TYPE "public"."BotType_old" RENAME TO "BotType"`);
    await queryRunner.query(`CREATE UNIQUE INDEX "users_tokenId_key" ON "users" ("tokenId") `);
    await queryRunner.query(
      `CREATE INDEX "dota_medals_rank_tier_idx" ON "dota_medals" ("rank_tier") `,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "dota_matches_match_id_key" ON "dota_matches" ("match_id") `,
    );
    await queryRunner.query(`CREATE UNIQUE INDEX "dota_heroes_id_key" ON "dota_heroes" ("id") `);
    await queryRunner.query(
      `CREATE INDEX "channels_dota_accounts_id_idx" ON "channels_dota_accounts" ("id") `,
    );
    await queryRunner.query(
      `CREATE INDEX "channels_commands_name_idx" ON "channels_commands" ("name") `,
    );
    await queryRunner.query(
      `CREATE INDEX "channels_commands_channelId_idx" ON "channels_commands" ("channelId") `,
    );
    await queryRunner.query(`CREATE UNIQUE INDEX "bots_tokenId_key" ON "bots" ("tokenId") `);
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ADD CONSTRAINT "users_viewed_notifications_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_viewed_notifications" ADD CONSTRAINT "users_viewed_notifications_notificationId_fkey" FOREIGN KEY ("notificationId") REFERENCES "notifications"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_files" ADD CONSTRAINT "users_files_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users" ADD CONSTRAINT "users_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications_messages" ADD CONSTRAINT "notifications_messages_notificationId_fkey" FOREIGN KEY ("notificationId") REFERENCES "notifications"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "notifications" ADD CONSTRAINT "notifications_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_results" ADD CONSTRAINT "dota_matches_results_match_id_fkey" FOREIGN KEY ("match_id") REFERENCES "dota_matches"("match_id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches_cards" ADD CONSTRAINT "dota_matches_cards_match_id_fkey" FOREIGN KEY ("match_id") REFERENCES "dota_matches"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "dota_matches" ADD CONSTRAINT "dota_matches_gameModeId_fkey" FOREIGN KEY ("gameModeId") REFERENCES "dota_game_modes"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "channels_dashboard_access_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "channels_dashboard_access_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE CASCADE ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ADD CONSTRAINT "channels_commands_responses_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE CASCADE ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ADD CONSTRAINT "channels_timers_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ADD CONSTRAINT "channels_permits_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_permits" ADD CONSTRAINT "channels_permits_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "channels_moderation_settings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ADD CONSTRAINT "channels_keywords_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_integrationId_fkey" FOREIGN KEY ("integrationId") REFERENCES "integrations"("id") ON DELETE CASCADE ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ADD CONSTRAINT "channels_greetings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "channels_dota_accounts_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_customvars" ADD CONSTRAINT "channels_customvars_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ADD CONSTRAINT "channels_commands_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" ADD CONSTRAINT "channels_id_fkey" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels" ADD CONSTRAINT "channels_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "bots" ADD CONSTRAINT "bots_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE SET NULL ON UPDATE CASCADE`,
    );
  }
}
