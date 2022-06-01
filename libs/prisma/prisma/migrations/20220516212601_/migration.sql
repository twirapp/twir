/*
  Warnings:

  - You are about to drop the column `interval` on the `timers` table. All the data in the column will be lost.
  - You are about to drop the column `triggerMessage` on the `timers` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "timers" DROP COLUMN "interval",
DROP COLUMN "triggerMessage",
ADD COLUMN     "messageInterval" INTEGER NOT NULL DEFAULT 0,
ADD COLUMN     "timeInterval" INTEGER NOT NULL DEFAULT 0;
