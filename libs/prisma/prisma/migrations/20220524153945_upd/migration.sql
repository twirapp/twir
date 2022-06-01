-- AlterTable
ALTER TABLE "channels_integrations" ALTER COLUMN "data" SET DATA TYPE JSONB;

-- AlterTable
ALTER TABLE "commands" ALTER COLUMN "aliases" SET DATA TYPE JSONB;

-- AlterTable
ALTER TABLE "timers" ALTER COLUMN "responses" SET DATA TYPE JSONB;

-- AlterTable
ALTER TABLE "users" ADD COLUMN     "testerId" TEXT;

-- CreateTable
CREATE TABLE "users_testers" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,

    CONSTRAINT "users_testers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ModerationSettings" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "links" JSONB NOT NULL DEFAULT '{}',
    "blacklist" JSONB NOT NULL DEFAULT '{}',
    "symbols" JSONB NOT NULL DEFAULT '{}',
    "longMessage" JSONB NOT NULL DEFAULT '{}',
    "caps" JSONB NOT NULL DEFAULT '{}',
    "emotes" JSONB NOT NULL DEFAULT '{}',

    CONSTRAINT "ModerationSettings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Permit" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,

    CONSTRAINT "Permit_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_testers_id_key" ON "users_testers"("id");

-- CreateIndex
CREATE UNIQUE INDEX "users_testers_userId_key" ON "users_testers"("userId");

-- CreateIndex
CREATE UNIQUE INDEX "ModerationSettings_channelId_key" ON "ModerationSettings"("channelId");

-- AddForeignKey
ALTER TABLE "users_testers" ADD CONSTRAINT "users_testers_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ModerationSettings" ADD CONSTRAINT "ModerationSettings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Permit" ADD CONSTRAINT "Permit_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Permit" ADD CONSTRAINT "Permit_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
