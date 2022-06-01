-- CreateEnum
CREATE TYPE "IntegrationService" AS ENUM ('LASTFM');

-- CreateTable
CREATE TABLE "Integration" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "service" "IntegrationService" NOT NULL,
    "accessToken" TEXT,
    "refreshToken" TEXT,
    "clientId" TEXT,
    "clientSecret" TEXT,
    "apiKey" TEXT,

    CONSTRAINT "Integration_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ChannelIntegration" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "channelId" TEXT NOT NULL,
    "integrationId" TEXT NOT NULL,
    "data" JSON NOT NULL,

    CONSTRAINT "ChannelIntegration_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "ChannelIntegration" ADD CONSTRAINT "ChannelIntegration_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ChannelIntegration" ADD CONSTRAINT "ChannelIntegration_integrationId_fkey" FOREIGN KEY ("integrationId") REFERENCES "Integration"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
