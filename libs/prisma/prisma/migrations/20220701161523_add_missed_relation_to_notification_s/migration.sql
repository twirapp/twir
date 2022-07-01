/*
  Warnings:

  - Added the required column `notificationId` to the `notifications_messages` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "notifications_messages" ADD COLUMN     "notificationId" TEXT NOT NULL;

-- AddForeignKey
ALTER TABLE "notifications_messages" ADD CONSTRAINT "notifications_messages_notificationId_fkey" FOREIGN KEY ("notificationId") REFERENCES "notifications"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
