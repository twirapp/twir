/*
  Warnings:

  - The `players` column on the `dota_matches` table would be dropped and recreated. This will lead to data loss if there is data in the column.

*/
-- AlterTable
ALTER TABLE "dota_matches" DROP COLUMN "players",
ADD COLUMN     "players" INTEGER[];
