import { MigrationInterface, QueryRunner } from 'typeorm';

export class removeTitleAndGame1683364542566 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`DELETE from "channels_commands" WHERE "defaultName" = $1`, ['title']);
    await queryRunner.query(`DELETE from "channels_commands" WHERE "defaultName" = $1`, ['game']);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {}
}
