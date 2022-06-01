/*
  Warnings:

  - You are about to drop the column `accessToken` on the `users` table. All the data in the column will be lost.
  - You are about to drop the column `refreshToken` on the `users` table. All the data in the column will be lost.
  - You are about to drop the `bots_tokens` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "bots_tokens" DROP CONSTRAINT "bots_tokens_botId_fkey";

-- AlterTable
ALTER TABLE "users" DROP COLUMN "accessToken",
DROP COLUMN "refreshToken";

-- DropTable
DROP TABLE "bots_tokens";

-- CreateTable
CREATE TABLE "tokens" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "accessToken" TEXT NOT NULL,
    "refreshToken" TEXT NOT NULL,
    "expiresIn" INTEGER NOT NULL,
    "obtainmentTimestamp" TIMESTAMP(3) NOT NULL,
    "botId" TEXT,
    "userId" TEXT,

    CONSTRAINT "tokens_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "tokens_botId_key" ON "tokens"("botId");

-- CreateIndex
CREATE UNIQUE INDEX "tokens_userId_key" ON "tokens"("userId");

-- AddForeignKey
ALTER TABLE "tokens" ADD CONSTRAINT "tokens_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "tokens" ADD CONSTRAINT "tokens_botId_fkey" FOREIGN KEY ("botId") REFERENCES "bots"("id") ON DELETE SET NULL ON UPDATE CASCADE;
