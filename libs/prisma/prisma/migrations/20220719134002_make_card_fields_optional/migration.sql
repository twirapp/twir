-- AlterTable
ALTER TABLE "dota_matches_cards" ALTER COLUMN "rank_tier" DROP NOT NULL,
ALTER COLUMN "leaderboard_rank" DROP NOT NULL;

-- AlterTable
ALTER TABLE "users_stats" ALTER COLUMN "watched" SET DEFAULT 0;
