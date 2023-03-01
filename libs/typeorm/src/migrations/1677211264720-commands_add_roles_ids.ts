import { MigrationInterface, QueryRunner } from "typeorm";

export class commandsAddRolesIds1677211264720 implements MigrationInterface {
    name = 'commandsAddRolesIds1677211264720'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "rolesIds" text array NOT NULL DEFAULT '{}'`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "rolesIds"`);
    }

}
