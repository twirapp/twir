import { MigrationInterface, QueryRunner } from 'typeorm';

export class removeWordsCounter1669396099043 implements MigrationInterface {
  name = 'removeWordsCounter1669396099043';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ADD "usages" integer NOT NULL DEFAULT '0'`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ALTER COLUMN "response" DROP NOT NULL`,
    );
    await queryRunner.query(`DROP TABLE "channels_words_counters"`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_keywords" ALTER COLUMN "response" SET NOT NULL`);
    await queryRunner.query(`ALTER TABLE "channels_keywords" DROP COLUMN "usages"`);
  }
}
