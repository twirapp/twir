import { MigrationInterface, QueryRunner } from "typeorm";

export class addEnabledToWordsCounters1669153151709 implements MigrationInterface {
    name = 'addEnabledToWordsCounters1669153151709'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_words_counters" DROP CONSTRAINT "FK_70b8ce8e96d7c5a6ef9da0a6de5"`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ADD "enabled" boolean NOT NULL DEFAULT true`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ADD CONSTRAINT "FK_ca0004881f53fcc413b24f8a738" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_words_counters" DROP CONSTRAINT "FK_ca0004881f53fcc413b24f8a738"`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" DROP COLUMN "enabled"`);
        await queryRunner.query(`ALTER TABLE "channels_words_counters" ADD CONSTRAINT "FK_70b8ce8e96d7c5a6ef9da0a6de5" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

}
