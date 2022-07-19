-- AlterTable
ALTER TABLE "dota_matches_results" ALTER COLUMN "radiant_win" DROP DEFAULT;

-- AlterTable
ALTER TABLE "users_stats" ALTER COLUMN "watched" SET DEFAULT 0;

-- CreateTable
CREATE TABLE "dota_matches_cards" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "match_id" TEXT NOT NULL,
    "account_id" TEXT NOT NULL,
    "rank_tier" INTEGER NOT NULL,
    "leaderboard_rank" INTEGER NOT NULL,

    CONSTRAINT "dota_matches_cards_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "dota_matches_cards_match_id_account_id_key" ON "dota_matches_cards"("match_id", "account_id");

-- AddForeignKey
ALTER TABLE "dota_matches_cards" ADD CONSTRAINT "dota_matches_cards_match_id_fkey" FOREIGN KEY ("match_id") REFERENCES "dota_matches"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
