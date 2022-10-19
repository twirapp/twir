import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelKeywordsIsReply1666189605642 implements MigrationInterface {
  name = 'addChannelKeywordsIsReply1666189605642';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_keywords" ADD "isReply" boolean NOT NULL DEFAULT false`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_keywords" DROP COLUMN "isReply"`);
  }
}
