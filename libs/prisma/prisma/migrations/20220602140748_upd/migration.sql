-- DropIndex
DROP INDEX "commands_channelId_name_idx";

-- CreateIndex
CREATE INDEX "commands_channelId_idx" ON "commands"("channelId");

-- CreateIndex
CREATE INDEX "commands_name_idx" ON "commands"("name");
