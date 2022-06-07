-- DropForeignKey
ALTER TABLE "bots" DROP CONSTRAINT "bots_tokenId_fkey";

-- DropForeignKey
ALTER TABLE "channels" DROP CONSTRAINT "channels_botId_fkey";

-- DropForeignKey
ALTER TABLE "channels" DROP CONSTRAINT "channels_id_fkey";

-- DropForeignKey
ALTER TABLE "channels_commands" DROP CONSTRAINT "channels_commands_channelId_fkey";

-- DropForeignKey
ALTER TABLE "channels_commands_usages" DROP CONSTRAINT "channels_commands_usages_userId_fkey";

-- DropForeignKey
ALTER TABLE "channels_greetings" DROP CONSTRAINT "channels_greetings_channelId_fkey";

-- DropForeignKey
ALTER TABLE "channels_integrations" DROP CONSTRAINT "channels_integrations_channelId_fkey";

-- DropForeignKey
ALTER TABLE "users" DROP CONSTRAINT "users_tokenId_fkey";

-- DropForeignKey
ALTER TABLE "users_stats" DROP CONSTRAINT "users_stats_channelId_fkey";

-- DropForeignKey
ALTER TABLE "users_stats" DROP CONSTRAINT "users_stats_userId_fkey";

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "users_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_id_fkey" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels" ADD CONSTRAINT "channels_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_integrations" ADD CONSTRAINT "channels_integrations_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "bots" ADD CONSTRAINT "bots_tokenId_fkey" FOREIGN KEY ("tokenId") REFERENCES "tokens"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_commands" ADD CONSTRAINT "channels_commands_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_commands_usages" ADD CONSTRAINT "channels_commands_usages_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "channels_greetings" ADD CONSTRAINT "channels_greetings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
