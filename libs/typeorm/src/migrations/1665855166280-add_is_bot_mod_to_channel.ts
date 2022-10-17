import { MigrationInterface, QueryRunner } from 'typeorm';

export class addIsBotModToChannel1665855166280 implements MigrationInterface {
  name = 'addIsBotModToChannel1665855166280';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels" ADD "isBotMod" boolean NOT NULL DEFAULT false`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels" DROP COLUMN "isBotMod"`);
  }
}
