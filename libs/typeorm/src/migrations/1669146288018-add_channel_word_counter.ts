import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelWordCounter1669146288018 implements MigrationInterface {
  name = 'addChannelWordCounter1669146288018';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "channels_words_counters" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "channelId" text NOT NULL, "phrase" text NOT NULL, "counter" integer NOT NULL, CONSTRAINT "PK_e10023d6d64448eff4a830f8dc8" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_words_counters" ADD CONSTRAINT "FK_70b8ce8e96d7c5a6ef9da0a6de5" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_words_counters" DROP CONSTRAINT "FK_70b8ce8e96d7c5a6ef9da0a6de5"`,
    );
    await queryRunner.query(`DROP TABLE "channels_words_counters"`);
  }
}
