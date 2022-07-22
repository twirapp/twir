-- AlterTable
ALTER TABLE "users_stats" ALTER COLUMN "watched" SET DEFAULT 0;

UPDATE "channels_timers" SET "timeInterval" = Floor("timeInterval" / 60);