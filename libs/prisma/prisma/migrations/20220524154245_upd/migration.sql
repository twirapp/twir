/*
  Warnings:

  - You are about to drop the `ModerationSettings` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Permit` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "ModerationSettings" DROP CONSTRAINT "ModerationSettings_channelId_fkey";

-- DropForeignKey
ALTER TABLE "Permit" DROP CONSTRAINT "Permit_channelId_fkey";

-- DropForeignKey
ALTER TABLE "Permit" DROP CONSTRAINT "Permit_userId_fkey";

-- DropTable
DROP TABLE "ModerationSettings";

-- DropTable
DROP TABLE "Permit";

-- CreateTable
CREATE TABLE "channels_moderation_settings" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "links" JSONB NOT NULL DEFAULT '{}',
    "blacklist" JSONB NOT NULL DEFAULT '{}',
    "symbols" JSONB NOT NULL DEFAULT '{}',
    "longMessage" JSONB NOT NULL DEFAULT '{}',
    "caps" JSONB NOT NULL DEFAULT '{}',
    "emotes" JSONB NOT NULL DEFAULT '{}',

    CONSTRAINT "channels_moderation_settings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "permits" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,

    CONSTRAINT "permits_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "channels_moderation_settings_channelId_key" ON "channels_moderation_settings"("channelId");

-- AddForeignKey
ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "channels_moderation_settings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permits" ADD CONSTRAINT "permits_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permits" ADD CONSTRAINT "permits_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
