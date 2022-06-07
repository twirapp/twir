-- CreateEnum
CREATE TYPE "IntegrationService" AS ENUM ('LASTFM', 'VK', 'FACEIT', 'SPOTIFY');

-- CreateEnum
CREATE TYPE "BotType" AS ENUM ('DEFAULT', 'CUSTOM');

-- CreateEnum
CREATE TYPE "CommandPermission" AS ENUM ('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER');

-- CreateEnum
CREATE TYPE "CooldownType" AS ENUM ('GLOBAL', 'PER_USER');

-- CreateEnum
CREATE TYPE "SettingsType" AS ENUM ('links', 'blacklists', 'symbols', 'longMessage', 'caps', 'emotes');

-- CreateEnum
CREATE TYPE "CustomVarType" AS ENUM ('SCRIPT', 'TEXT');

-- CreateTable
CREATE TABLE "users" (
    "id" TEXT NOT NULL,
    "tokenId" TEXT,
    "isTester" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "users_stats" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "userId" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,
    "messages" INTEGER NOT NULL DEFAULT 0,
    "watched" BIGINT NOT NULL DEFAULT 0,

    CONSTRAINT "users_stats_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels" (
    "id" TEXT NOT NULL,
    "isEnabled" BOOLEAN NOT NULL DEFAULT true,
    "isTwitchBanned" BOOLEAN NOT NULL DEFAULT false,
    "isBanned" BOOLEAN NOT NULL DEFAULT false,
    "botId" TEXT NOT NULL,

    CONSTRAINT "channels_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_dashboard_access" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,

    CONSTRAINT "channels_dashboard_access_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "integrations" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "service" "IntegrationService" NOT NULL,
    "accessToken" TEXT,
    "refreshToken" TEXT,
    "clientId" TEXT,
    "clientSecret" TEXT,
    "apiKey" TEXT,
    "redirectUrl" TEXT,

    CONSTRAINT "integrations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_integrations" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "channelId" TEXT NOT NULL,
    "integrationId" TEXT NOT NULL,
    "accessToken" TEXT,
    "refreshToken" TEXT,
    "clientId" TEXT,
    "clientSecret" TEXT,
    "apiKey" TEXT,
    "data" JSONB,

    CONSTRAINT "channels_integrations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "bots" (
    "id" TEXT NOT NULL,
    "type" "BotType" NOT NULL,
    "tokenId" TEXT,

    CONSTRAINT "bots_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_commands" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "cooldown" INTEGER DEFAULT 0,
    "cooldownType" "CooldownType" NOT NULL DEFAULT E'GLOBAL',
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "aliases" JSONB DEFAULT '[]',
    "description" TEXT,
    "visible" BOOLEAN NOT NULL DEFAULT true,
    "channelId" TEXT NOT NULL,
    "permission" "CommandPermission" NOT NULL,
    "default" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "channels_commands_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_commands_responses" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "text" TEXT,
    "commandId" TEXT NOT NULL,

    CONSTRAINT "channels_commands_responses_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_commands_usages" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "userId" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,
    "commandId" TEXT NOT NULL,

    CONSTRAINT "channels_commands_usages_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_timers" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "responses" JSONB NOT NULL DEFAULT '[]',
    "last" INTEGER NOT NULL DEFAULT 0,
    "timeInterval" INTEGER NOT NULL DEFAULT 0,
    "messageInterval" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "channels_timers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "tokens" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "accessToken" TEXT NOT NULL,
    "refreshToken" TEXT NOT NULL,
    "expiresIn" INTEGER NOT NULL,
    "obtainmentTimestamp" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "tokens_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_greetings" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "text" TEXT NOT NULL,

    CONSTRAINT "channels_greetings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_moderation_settings" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "type" "SettingsType" NOT NULL,
    "channelId" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "subscribers" BOOLEAN NOT NULL DEFAULT false,
    "vips" BOOLEAN NOT NULL DEFAULT false,
    "banTime" INTEGER NOT NULL,
    "banMessage" TEXT NOT NULL,
    "warningMessage" TEXT NOT NULL,
    "checkClips" BOOLEAN,
    "triggerLength" INTEGER,
    "maxPercentage" INTEGER,
    "blackListSentences" JSONB,

    CONSTRAINT "channels_moderation_settings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_permits" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,

    CONSTRAINT "channels_permits_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_keywords" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "response" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "cooldown" INTEGER DEFAULT 0,

    CONSTRAINT "channels_keywords_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels_customvars" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "description" TEXT,
    "type" "CustomVarType" NOT NULL,
    "evalValue" TEXT,
    "response" TEXT,
    "channelId" TEXT NOT NULL,

    CONSTRAINT "channels_customvars_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_tokenId_key" ON "users"("tokenId");

-- CreateIndex
CREATE UNIQUE INDEX "users_stats_userId_channelId_key" ON "users_stats"("userId", "channelId");

-- CreateIndex
CREATE UNIQUE INDEX "bots_tokenId_key" ON "bots"("tokenId");

-- CreateIndex
CREATE INDEX "channels_commands_channelId_idx" ON "channels_commands"("channelId");

-- CreateIndex
CREATE INDEX "channels_commands_name_idx" ON "channels_commands"("name");

-- CreateIndex
CREATE UNIQUE INDEX "channels_moderation_settings_channelId_key" ON "channels_moderation_settings"("channelId");

-- CreateIndex
CREATE UNIQUE INDEX "channels_moderation_settings_channelId_type_key" ON "channels_moderation_settings"("channelId", "type");

-- CreateIndex
CREATE UNIQUE INDEX "channels_keywords_channelId_text_key" ON "channels_keywords"("channelId", "text");

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "users_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_id_fkey" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "channels_dashboard_access_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "channels_dashboard_access_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_integrationId_fkey" FOREIGN KEY ("integrationId") REFERENCES "integrations"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "bots" ADD CONSTRAINT "bots_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_commands" ADD CONSTRAINT "channels_commands_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_commands_responses" ADD CONSTRAINT "channels_commands_responses_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_timers" ADD CONSTRAINT "channels_timers_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_greetings" ADD CONSTRAINT "channels_greetings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "channels_moderation_settings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_permits" ADD CONSTRAINT "channels_permits_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_permits" ADD CONSTRAINT "channels_permits_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_keywords" ADD CONSTRAINT "channels_keywords_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_customvars" ADD CONSTRAINT "channels_customvars_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
