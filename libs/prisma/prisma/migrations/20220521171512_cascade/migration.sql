-- DropForeignKey
ALTER TABLE "ChannelIntegration" DROP CONSTRAINT "ChannelIntegration_channelId_fkey";

-- DropForeignKey
ALTER TABLE "ChannelIntegration" DROP CONSTRAINT "ChannelIntegration_integrationId_fkey";

-- DropForeignKey
ALTER TABLE "bots_tokens" DROP CONSTRAINT "bots_tokens_botId_fkey";

-- DropForeignKey
ALTER TABLE "channels" DROP CONSTRAINT "channels_botId_fkey";

-- DropForeignKey
ALTER TABLE "channels" DROP CONSTRAINT "channels_id_fkey";

-- DropForeignKey
ALTER TABLE "commands" DROP CONSTRAINT "commands_channelId_fkey";

-- DropForeignKey
ALTER TABLE "commands_responses" DROP CONSTRAINT "commands_responses_commandId_fkey";

-- DropForeignKey
ALTER TABLE "commands_usages" DROP CONSTRAINT "commands_usages_commandId_fkey";

-- DropForeignKey
ALTER TABLE "commands_usages" DROP CONSTRAINT "commands_usages_userId_fkey";

-- DropForeignKey
ALTER TABLE "greetings" DROP CONSTRAINT "greetings_channelId_fkey";

-- DropForeignKey
ALTER TABLE "users_stats" DROP CONSTRAINT "users_stats_channelId_fkey";

-- DropForeignKey
ALTER TABLE "users_stats" DROP CONSTRAINT "users_stats_userId_fkey";

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_id_fkey" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ChannelIntegration" ADD CONSTRAINT "ChannelIntegration_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ChannelIntegration" ADD CONSTRAINT "ChannelIntegration_integrationId_fkey" FOREIGN KEY ("integrationId") REFERENCES "Integration"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands" ADD CONSTRAINT "commands_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands_responses" ADD CONSTRAINT "commands_responses_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "commands"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands_usages" ADD CONSTRAINT "commands_usages_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands_usages" ADD CONSTRAINT "commands_usages_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "commands"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "bots_tokens" ADD CONSTRAINT "bots_tokens_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "greetings" ADD CONSTRAINT "greetings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE CASCADE;
