-- CreateEnum
CREATE TYPE "CommandModule" AS ENUM ('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION');

-- AlterTable
ALTER TABLE "channels_commands" ADD COLUMN     "module" "CommandModule" NOT NULL DEFAULT 'CUSTOM';

-- AlterTable
ALTER TABLE "users_stats" ALTER COLUMN "watched" SET DEFAULT 0;

UPDATE channels_commands
SET "module" = 'DOTA'
WHERE "defaultName" IN ('dota addacc', 'dota delacc', 'np', 'wl', 'dota listacc', 'lg', 'gm');

UPDATE channels_commands
SET "module" = 'CHANNEL'
WHERE "defaultName" IN ('title set', 'game set');

UPDATE channels_commands
SET "module" = 'MODERATION'
WHERE "defaultName" IN ('permit');