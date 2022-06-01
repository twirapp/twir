-- CreateEnum
CREATE TYPE "BotType" AS ENUM ('DEFAULT', 'CUSTOM');

-- CreateEnum
CREATE TYPE "CommandPermission" AS ENUM ('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER');

-- CreateTable
CREATE TABLE "users" (
    "id" TEXT NOT NULL,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "channels" (
    "id" TEXT NOT NULL,
    "isEnabled" BOOLEAN NOT NULL DEFAULT true,
    "botId" TEXT NOT NULL,

    CONSTRAINT "channels_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "bots" (
    "id" TEXT NOT NULL,
    "type" "BotType" NOT NULL,

    CONSTRAINT "bots_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "commands" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "cooldown" INTEGER DEFAULT 0,
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "aliases" JSON DEFAULT '[]',
    "description" TEXT,
    "visible" BOOLEAN NOT NULL DEFAULT true,
    "channelId" TEXT NOT NULL,
    "permission" "CommandPermission" NOT NULL,

    CONSTRAINT "commands_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "commands_responses" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "text" TEXT,
    "commandId" TEXT NOT NULL,

    CONSTRAINT "commands_responses_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "timers" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "interval" INTEGER NOT NULL DEFAULT 0,
    "responses" JSON NOT NULL DEFAULT '[]',
    "last" INTEGER,
    "triggerTimeStamp" BIGINT DEFAULT 0,
    "triggerMessage" INTEGER DEFAULT 0,

    CONSTRAINT "timers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "bots_tokens" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "botId" TEXT NOT NULL,
    "accessToken" TEXT NOT NULL,
    "refreshToken" TEXT NOT NULL,
    "expiresIn" INTEGER NOT NULL,
    "obtainmentTimestamp" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "bots_tokens_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "greetings" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "userId" INTEGER,
    "username" TEXT,

    CONSTRAINT "greetings_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "bots_tokens_botId_key" ON "bots_tokens"("botId");

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_id_fkey" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands" ADD CONSTRAINT "commands_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands_responses" ADD CONSTRAINT "commands_responses_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "commands"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "timers" ADD CONSTRAINT "timers_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "bots_tokens" ADD CONSTRAINT "bots_tokens_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
