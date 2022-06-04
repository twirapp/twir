-- DropIndex
DROP INDEX "commands_channelId_idx";

-- AlterTable
ALTER TABLE "commands" ADD COLUMN     "default" BOOLEAN NOT NULL DEFAULT false;

-- CreateIndex
CREATE INDEX "commands_channelId_name_idx" ON "commands"("channelId", "name");
