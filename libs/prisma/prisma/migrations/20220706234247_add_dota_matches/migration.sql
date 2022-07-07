-- CreateTable
CREATE TABLE "dota_matches" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "startedAt" TIMESTAMP(3) NOT NULL,
    "lobby_type" INTEGER NOT NULL,
    "gameModeId" TEXT NOT NULL,
    "players" JSONB[],
    "weekend_tourney_bracket_round" TEXT,
    "weekend_tourney_skill_level" TEXT,
    "match_id" TEXT NOT NULL,

    CONSTRAINT "dota_matches_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "dota_matches" ADD CONSTRAINT "dota_matches_gameModeId_fkey" FOREIGN KEY ("gameModeId") REFERENCES "dota_game_modes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
