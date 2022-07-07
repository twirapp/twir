/*
  Warnings:

  - The primary key for the `dota_heroes` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - A unique constraint covering the columns `[id]` on the table `dota_heroes` will be added. If there are existing duplicate values, this will fail.
  - Changed the type of `id` on the `dota_heroes` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- AlterTable
ALTER TABLE "dota_heroes" DROP CONSTRAINT "dota_heroes_pkey",
DROP COLUMN "id",
ADD COLUMN     "id" INTEGER NOT NULL,
ADD CONSTRAINT "dota_heroes_pkey" PRIMARY KEY ("id");

-- CreateIndex
CREATE UNIQUE INDEX "dota_heroes_id_key" ON "dota_heroes"("id");
