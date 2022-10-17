import { MigrationInterface, QueryRunner } from 'typeorm';

export class addKeywordsCooldownExpireAt1665875488674 implements MigrationInterface {
  name = 'addKeywordsCooldownExpireAt1665875488674';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_keywords" ADD "cooldownExpireAt" TIMESTAMP`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_keywords" DROP COLUMN "cooldownExpireAt"`);
  }
}
