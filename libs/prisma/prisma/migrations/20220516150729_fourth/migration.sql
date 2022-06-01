/*
  Warnings:

  - Added the required column `text` to the `greetings` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "greetings" ADD COLUMN     "text" TEXT NOT NULL;
