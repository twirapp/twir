/*
  Warnings:

  - A unique constraint covering the columns `[match_id]` on the table `dota_matches` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "dota_matches_match_id_key" ON "dota_matches"("match_id");
