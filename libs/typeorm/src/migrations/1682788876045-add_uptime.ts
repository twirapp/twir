import { MigrationInterface, QueryRunner } from 'typeorm';

export class addUptime1682788876045 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `DELETE from "channels_commands" WHERE "name" = $1 AND "module" = 'CUSTOM'`,
      ['uptime'],
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {}
}
