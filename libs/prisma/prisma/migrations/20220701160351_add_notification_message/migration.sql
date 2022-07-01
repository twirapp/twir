/*
  Warnings:

  - You are about to drop the column `text` on the `notifications` table. All the data in the column will be lost.
  - You are about to drop the column `title` on the `notifications` table. All the data in the column will be lost.

*/
-- CreateEnum
CREATE TYPE "LangCode" AS ENUM ('RU', 'GB');

-- AlterTable
ALTER TABLE "notifications" DROP COLUMN "text",
DROP COLUMN "title";

-- CreateTable
CREATE TABLE "notifications_messages" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "text" TEXT NOT NULL,
    "title" TEXT,
    "langCode" "LangCode" NOT NULL,

    CONSTRAINT "notifications_messages_pkey" PRIMARY KEY ("id")
);
