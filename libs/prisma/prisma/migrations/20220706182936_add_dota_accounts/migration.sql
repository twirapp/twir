-- AlterEnum
ALTER TYPE "IntegrationService" ADD VALUE 'DOTA';

-- CreateTable
CREATE TABLE "channels_dota_accounts" (
    "id" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,

    CONSTRAINT "channels_dota_accounts_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "channels_dota_accounts_id_channelId_key" ON "channels_dota_accounts"("id", "channelId");

-- AddForeignKey
ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "channels_dota_accounts_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
