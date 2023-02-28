import { MigrationInterface, QueryRunner } from "typeorm";

export class removeUnnecessarySystemColumnFromRoles1677598576059 implements MigrationInterface {
    name = 'removeUnnecessarySystemColumnFromRoles1677598576059'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles" DROP COLUMN "system"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles" ADD "system" boolean NOT NULL DEFAULT false`);
    }

}
