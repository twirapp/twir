-- CreateTable
CREATE TABLE "channels_dota_accounts" (
    "id" TEXT NOT NULL,
    "channelId" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "dota_heroes" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "dota_heroes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "dota_game_modes" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "dota_game_modes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "dota_medals" (
    "rank_tier" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "dota_medals_pkey" PRIMARY KEY ("rank_tier")
);

-- CreateTable
CREATE TABLE "dota_matches" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "startedAt" TIMESTAMP(3) NOT NULL,
    "lobby_type" INTEGER NOT NULL,
    "gameModeId" INTEGER NOT NULL,
    "players" INTEGER[],
    "players_heroes" INTEGER[],
    "weekend_tourney_bracket_round" TEXT,
    "weekend_tourney_skill_level" TEXT,
    "match_id" TEXT NOT NULL,
    "avarage_mmr" INTEGER NOT NULL,
    "lobbyId" TEXT NOT NULL,
    "finished" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "dota_matches_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "channels_dota_accounts_id_idx" ON "channels_dota_accounts"("id");

-- CreateIndex
CREATE UNIQUE INDEX "channels_dota_accounts_id_channelId_key" ON "channels_dota_accounts"("id", "channelId");

-- CreateIndex
CREATE UNIQUE INDEX "dota_heroes_id_key" ON "dota_heroes"("id");

-- CreateIndex
CREATE INDEX "dota_medals_rank_tier_idx" ON "dota_medals"("rank_tier");

-- CreateIndex
CREATE UNIQUE INDEX "dota_matches_match_id_key" ON "dota_matches"("match_id");

-- AddForeignKey
ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "channels_dota_accounts_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "dota_matches" ADD CONSTRAINT "dota_matches_gameModeId_fkey" FOREIGN KEY ("gameModeId") REFERENCES "dota_game_modes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
