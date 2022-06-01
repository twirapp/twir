/*
  Warnings:

  - You are about to drop the column `blacklist` on the `channels_moderation_settings` table. All the data in the column will be lost.
  - You are about to drop the column `caps` on the `channels_moderation_settings` table. All the data in the column will be lost.
  - You are about to drop the column `emotes` on the `channels_moderation_settings` table. All the data in the column will be lost.
  - You are about to drop the column `links` on the `channels_moderation_settings` table. All the data in the column will be lost.
  - You are about to drop the column `longMessage` on the `channels_moderation_settings` table. All the data in the column will be lost.
  - You are about to drop the column `symbols` on the `channels_moderation_settings` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[channelId,type]` on the table `channels_moderation_settings` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `banMessage` to the `channels_moderation_settings` table without a default value. This is not possible if the table is not empty.
  - Added the required column `banTime` to the `channels_moderation_settings` table without a default value. This is not possible if the table is not empty.
  - Added the required column `type` to the `channels_moderation_settings` table without a default value. This is not possible if the table is not empty.
  - Added the required column `warningMessage` to the `channels_moderation_settings` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "SettingsType" AS ENUM ('links', 'blacklists', 'symbols', 'longMessage', 'caps', 'emotes');

-- AlterTable
ALTER TABLE "channels_moderation_settings" DROP COLUMN "blacklist",
DROP COLUMN "caps",
DROP COLUMN "emotes",
DROP COLUMN "links",
DROP COLUMN "longMessage",
DROP COLUMN "symbols",
ADD COLUMN     "banMessage" TEXT NOT NULL,
ADD COLUMN     "banTime" INTEGER NOT NULL,
ADD COLUMN     "blackListSentences" JSONB,
ADD COLUMN     "checkClips" BOOLEAN,
ADD COLUMN     "maxPercentage" INTEGER,
ADD COLUMN     "subscribers" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "triggerLength" INTEGER,
ADD COLUMN     "type" "SettingsType" NOT NULL,
ADD COLUMN     "vips" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "warningMessage" TEXT NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "channels_moderation_settings_channelId_type_key" ON "channels_moderation_settings"("channelId", "type");
