-- CreateTable
CREATE TABLE "channels_dota_accounts" (
    "id" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,

    CONSTRAINT "channels_dota_accounts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "dota_heroes" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "dota_heroes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "dota_game_modes" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "dota_game_modes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "dota_medals" (
    "rank_tier" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "dota_medals_pkey" PRIMARY KEY ("rank_tier")
);

-- CreateIndex
CREATE UNIQUE INDEX "channels_dota_accounts_id_channelId_key" ON "channels_dota_accounts"("id", "channelId");

-- CreateIndex
CREATE INDEX "dota_medals_rank_tier_idx" ON "dota_medals"("rank_tier");

-- AddForeignKey
ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "channels_dota_accounts_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
