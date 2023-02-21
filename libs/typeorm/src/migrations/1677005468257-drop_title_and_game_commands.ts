import { MigrationInterface, QueryRunner } from 'typeorm';

export class dropTitleAndGameCommands1677005468257 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`DELETE FROM "channels_commands" WHERE "defaultName" = $1`, ['title set']);
        await queryRunner.query(`DELETE FROM "channels_commands" WHERE "defaultName" = $1`, ['game set']);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
