import { MigrationInterface, QueryRunner } from "typeorm";

export class eventOperationsEnabled1686603414357 implements MigrationInterface {
    name = 'eventOperationsEnabled1686603414357'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ADD "enabled" boolean NOT NULL DEFAULT true`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations" DROP COLUMN "enabled"`);
    }

}
