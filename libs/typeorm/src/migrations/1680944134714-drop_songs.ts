import { MigrationInterface, QueryRunner } from "typeorm"

export class dropSongs1680944134714 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(
          `DELETE from "channels_commands" WHERE "name" = $1 AND "module" = 'CUSTOM'`,
          ['song'],
        );
        await queryRunner.query(
          `DELETE from "channels_commands" WHERE "name" = $1 AND "module" = 'CUSTOM'`,
          ['currentsong'],
        );
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
