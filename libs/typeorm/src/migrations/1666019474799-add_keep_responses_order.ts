import { MigrationInterface, QueryRunner } from 'typeorm';

export class addKeepResponsesOrder1666019474799 implements MigrationInterface {
  name = 'addKeepResponsesOrder1666019474799';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ADD "keepResponsesOrder" boolean NOT NULL DEFAULT true`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "keepResponsesOrder"`);
  }
}
