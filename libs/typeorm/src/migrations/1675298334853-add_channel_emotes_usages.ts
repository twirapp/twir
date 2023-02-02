import { MigrationInterface, QueryRunner } from "typeorm";

export class addChannelEmotesUsages1675298334853 implements MigrationInterface {
    name = 'addChannelEmotesUsages1675298334853'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "channels_emotes_usages" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "channelId" text NOT NULL, "userId" text NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "emote" character varying NOT NULL, CONSTRAINT "PK_cbaa7cf66062bb2a4d7926826f8" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" ADD CONSTRAINT "FK_309ade49a31238d00065fc7c32e" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" ADD CONSTRAINT "FK_a736d0930dd0d26bd5f52bb0cf0" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" DROP CONSTRAINT "FK_a736d0930dd0d26bd5f52bb0cf0"`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" DROP CONSTRAINT "FK_309ade49a31238d00065fc7c32e"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`DROP TABLE "channels_emotes_usages"`);
    }

}
