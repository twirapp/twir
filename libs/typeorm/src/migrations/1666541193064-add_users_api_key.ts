import { MigrationInterface, QueryRunner } from 'typeorm';

export class addUsersApiKey1666541193064 implements MigrationInterface {
  name = 'addUsersApiKey1666541193064';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "users" ADD "apiKey" uuid NOT NULL DEFAULT gen_random_uuid()`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "users" DROP COLUMN "apiKey"`);
  }
}
