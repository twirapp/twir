/*
  Warnings:

  - Added the required column `avarage_mmr` to the `dota_matches` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "dota_matches" ADD COLUMN     "avarage_mmr" INTEGER NOT NULL;
