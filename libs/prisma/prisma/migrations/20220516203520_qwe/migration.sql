/*
  Warnings:

  - You are about to drop the column `triggerTimeStamp` on the `timers` table. All the data in the column will be lost.
  - Made the column `last` on table `timers` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE "timers" DROP COLUMN "triggerTimeStamp",
ALTER COLUMN "last" SET NOT NULL,
ALTER COLUMN "last" SET DEFAULT 0;
