/*
  Warnings:

  - You are about to drop the `ChannelIntegration` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Integration` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "ChannelIntegration" DROP CONSTRAINT "ChannelIntegration_channelId_fkey";

-- DropForeignKey
ALTER TABLE "ChannelIntegration" DROP CONSTRAINT "ChannelIntegration_integrationId_fkey";

-- DropTable
DROP TABLE "ChannelIntegration";

-- DropTable
DROP TABLE "Integration";

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
    "data" JSON,

    CONSTRAINT "channels_integrations_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_integrationId_fkey" FOREIGN KEY ("integrationId") REFERENCES "integrations"("id") ON DELETE CASCADE ON UPDATE CASCADE;
