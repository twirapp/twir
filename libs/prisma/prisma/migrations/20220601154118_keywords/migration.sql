-- CreateTable
CREATE TABLE "keywords" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "response" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "cooldown" INTEGER DEFAULT 0,

    CONSTRAINT "keywords_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "keywords" ADD CONSTRAINT "keywords_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
