/*
  Warnings:

  - The primary key for the `dota_game_modes` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - Changed the type of `id` on the `dota_game_modes` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.
  - Changed the type of `gameModeId` on the `dota_matches` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- DropForeignKey
ALTER TABLE "dota_matches" DROP CONSTRAINT "dota_matches_gameModeId_fkey";

-- AlterTable
ALTER TABLE "dota_game_modes" DROP CONSTRAINT "dota_game_modes_pkey",
DROP COLUMN "id",
ADD COLUMN     "id" INTEGER NOT NULL,
ADD CONSTRAINT "dota_game_modes_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "dota_matches" DROP COLUMN "gameModeId",
ADD COLUMN     "gameModeId" INTEGER NOT NULL;

-- AddForeignKey
ALTER TABLE "dota_matches" ADD CONSTRAINT "dota_matches_gameModeId_fkey" FOREIGN KEY ("gameModeId") REFERENCES "dota_game_modes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
