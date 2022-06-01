-- DropIndex
DROP INDEX "users_testers_id_key";

-- AlterTable
ALTER TABLE "commands" ALTER COLUMN "aliases" SET DEFAULT '[]';

-- AlterTable
ALTER TABLE "timers" ALTER COLUMN "responses" SET DEFAULT '[]';

-- AlterTable
ALTER TABLE "users_testers" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
