import { MigrationInterface, QueryRunner } from 'typeorm';

export class addOrderToCommandsResponses1666019651002 implements MigrationInterface {
  name = 'addOrderToCommandsResponses1666019651002';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_commands_responses" ADD "order" integer NOT NULL DEFAULT '0'`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_commands_responses" DROP COLUMN "order"`);
  }
}
