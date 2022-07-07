/*
  Warnings:

  - Changed the type of `players` on the `dota_matches` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- AlterTable
ALTER TABLE "dota_matches" DROP COLUMN "players",
ADD COLUMN     "players" JSONB NOT NULL;
