import { MigrationInterface, QueryRunner } from 'typeorm';

export class addModerationWarnings1665966014918 implements MigrationInterface {
  name = 'addModerationWarnings1665966014918';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "channels_moderation_warnings" ("id" text NOT NULL DEFAULT gen_random_uuid(), "channelId" text NOT NULL, "userId" text NOT NULL, "reason" text NOT NULL, CONSTRAINT "PK_311c956bfd3e98c159ffb78740d" PRIMARY KEY ("id"))`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`DROP TABLE "channels_moderation_warnings"`);
  }
}
