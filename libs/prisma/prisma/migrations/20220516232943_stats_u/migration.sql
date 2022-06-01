/*
  Warnings:

  - A unique constraint covering the columns `[userId,channelId]` on the table `users_stats` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "users_stats_userId_channelId_key" ON "users_stats"("userId", "channelId");
