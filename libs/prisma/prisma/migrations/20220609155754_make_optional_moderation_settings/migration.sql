-- AlterTable
ALTER TABLE "channels_moderation_settings" ALTER COLUMN "banTime" SET DEFAULT 600,
ALTER COLUMN "banMessage" DROP NOT NULL,
ALTER COLUMN "warningMessage" DROP NOT NULL,
ALTER COLUMN "checkClips" SET DEFAULT false,
ALTER COLUMN "triggerLength" SET DEFAULT 300,
ALTER COLUMN "maxPercentage" SET DEFAULT 50,
ALTER COLUMN "blackListSentences" SET DEFAULT '[]';
