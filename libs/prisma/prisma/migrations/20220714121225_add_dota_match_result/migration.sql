-- AlterTable
ALTER TABLE "users_stats" ALTER COLUMN "watched" SET DEFAULT 0;

-- CreateTable
CREATE TABLE "dota_matches_results" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "match_id" TEXT NOT NULL,
    "players" JSONB NOT NULL,
    "radiant_win" BOOLEAN NOT NULL DEFAULT false,
    "game_mode" INTEGER NOT NULL,

    CONSTRAINT "dota_matches_results_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "dota_matches_results_match_id_key" ON "dota_matches_results"("match_id");

-- AddForeignKey
ALTER TABLE "dota_matches_results" ADD CONSTRAINT "dota_matches_results_match_id_fkey" FOREIGN KEY ("match_id") REFERENCES "dota_matches"("match_id") ON DELETE RESTRICT ON UPDATE CASCADE;
