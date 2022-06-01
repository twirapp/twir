/*
  Warnings:

  - You are about to drop the column `username` on the `greetings` table. All the data in the column will be lost.
  - Made the column `userId` on table `greetings` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE "greetings" DROP COLUMN "username",
ALTER COLUMN "userId" SET NOT NULL;
