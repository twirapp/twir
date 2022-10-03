import { MigrationInterface, QueryRunner } from 'typeorm';

export class addIsReplyColumnToCommands1664809110625 implements MigrationInterface {
  name = 'addIsReplyColumnToCommands1664809110625';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ADD "is_reply" boolean NOT NULL DEFAULT true`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "is_reply"`);
  }
}
