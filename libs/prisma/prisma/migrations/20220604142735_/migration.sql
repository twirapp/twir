/*
  Warnings:

  - A unique constraint covering the columns `[channelId,text]` on the table `keywords` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "keywords_channelId_text_key" ON "keywords"("channelId", "text");
