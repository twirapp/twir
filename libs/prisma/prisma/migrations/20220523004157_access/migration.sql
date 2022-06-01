-- CreateTable
CREATE TABLE "dashboard_access" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,

    CONSTRAINT "dashboard_access_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "dashboard_access" ADD CONSTRAINT "dashboard_access_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "dashboard_access" ADD CONSTRAINT "dashboard_access_channelId_fkey" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
