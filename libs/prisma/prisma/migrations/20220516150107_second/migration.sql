/*
  Warnings:

  - Added the required column `channelId` to the `greetings` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "greetings" ADD COLUMN     "channelId" TEXT NOT NULL,
ALTER COLUMN "userId" SET DATA TYPE TEXT;

-- CreateTable
CREATE TABLE "commands_usages" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "userId" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,
    "commandId" TEXT NOT NULL,

    CONSTRAINT "commands_usages_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "commands_channelId_idx" ON "commands"("channelId");

-- AddForeignKey
ALTER TABLE "commands_usages" ADD CONSTRAINT "commands_usages_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "commands_usages" ADD CONSTRAINT "commands_usages_commandId_fkey" FOREIGN KEY ("commandId") REFERENCES "commands"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "greetings" ADD CONSTRAINT "greetings_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
