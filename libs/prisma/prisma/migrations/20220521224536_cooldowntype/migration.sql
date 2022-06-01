-- CreateEnum
CREATE TYPE "CooldownType" AS ENUM ('GLOBAL', 'PER_USER');

-- AlterTable
ALTER TABLE "commands" ADD COLUMN     "cooldownType" "CooldownType" NOT NULL DEFAULT E'GLOBAL';
