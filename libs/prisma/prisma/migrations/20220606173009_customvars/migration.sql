-- CreateEnum
CREATE TYPE "CustomVarType" AS ENUM ('SCRIPT', 'TEXT');

-- CreateTable
CREATE TABLE "CustomVar" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "description" TEXT,
    "type" "CustomVarType" NOT NULL,
    "evalValue" TEXT,
    "response" TEXT,

    CONSTRAINT "CustomVar_pkey" PRIMARY KEY ("id")
);
