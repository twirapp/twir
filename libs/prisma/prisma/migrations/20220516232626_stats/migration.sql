-- CreateTable
CREATE TABLE "users_stats" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,
    "messages" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "users_stats_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
