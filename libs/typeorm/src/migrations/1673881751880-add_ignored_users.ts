import { MigrationInterface, QueryRunner } from "typeorm";

export class addIgnoredUsers1673881751880 implements MigrationInterface {
    name = 'addIgnoredUsers1673881751880'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "users_ignored" ("id" text NOT NULL, "login" text, "displayName" text, CONSTRAINT "PK_ac130ac65d82c39fee96465fb8d" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE INDEX "IDX_ff17ac70a0b164274148218375" ON "users_ignored" ("login") `);
        await queryRunner.query(`CREATE INDEX "IDX_74a54fd2a9d27143a7ed6bb979" ON "users_ignored" ("displayName") `);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`DROP INDEX "public"."IDX_74a54fd2a9d27143a7ed6bb979"`);
        await queryRunner.query(`DROP INDEX "public"."IDX_ff17ac70a0b164274148218375"`);
        await queryRunner.query(`DROP TABLE "users_ignored"`);
    }

}
