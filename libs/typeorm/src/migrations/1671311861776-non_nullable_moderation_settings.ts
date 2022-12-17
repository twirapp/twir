import { MigrationInterface, QueryRunner } from 'typeorm';

export class nonNullableModerationSettings1671311861776 implements MigrationInterface {
  name = 'nonNullableModerationSettings1671311861776';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "banMessage" SET DEFAULT ''`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "warningMessage" SET DEFAULT ''`,
    );

    await queryRunner.query(
      `UPDATE "channels_moderation_settings" SET "banMessage"=$1 WHERE "banMessage" IS NULL`,
      [''],
    );

    await queryRunner.query(
      `UPDATE "channels_moderation_settings" SET "warningMessage"=$1 WHERE "warningMessage" IS NULL`,
      [''],
    );

    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "banMessage" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "warningMessage" SET NOT NULL`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "warningMessage" DROP DEFAULT`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "banMessage" DROP DEFAULT`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "warningMessage" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "banMessage" DROP NOT NULL`,
    );
  }
}
