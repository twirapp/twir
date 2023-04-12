import { MigrationInterface, QueryRunner } from 'typeorm';

export class addAccountage1681298402900 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(
            `DELETE from "channels_commands" WHERE "name" = $1 AND "module" = 'CUSTOM'`,
            ['age'],
          );
          await queryRunner.query(
            `DELETE from "channels_commands" WHERE "name" = $1 AND "module" = 'CUSTOM'`,
            ['accountage'],
          );
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
