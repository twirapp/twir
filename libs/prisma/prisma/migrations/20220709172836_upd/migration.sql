/*
  Warnings:

  - The primary key for the `channels_dota_accounts` table will be changed. If it partially fails, the table could be left without primary key constraint.

*/
-- AlterTable
ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "channels_dota_accounts_pkey";

-- CreateIndex
CREATE INDEX "channels_dota_accounts_id_idx" ON "channels_dota_accounts"("id");
