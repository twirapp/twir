/*
  Warnings:

  - You are about to drop the `CustomVar` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropTable
DROP TABLE "CustomVar";

-- CreateTable
CREATE TABLE "channels_customvars" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "type" "CustomVarType" NOT NULL,
    "evalValue" TEXT,
    "response" TEXT,

    CONSTRAINT "channels_customvars_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "channels_customvars" ADD CONSTRAINT "channels_customvars_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
