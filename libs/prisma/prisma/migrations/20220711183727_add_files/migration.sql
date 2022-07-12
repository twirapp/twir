-- CreateTable
CREATE TABLE "users_files" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "size" INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    "userId" TEXT,

    CONSTRAINT "users_files_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "users_files" ADD CONSTRAINT "users_files_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;
