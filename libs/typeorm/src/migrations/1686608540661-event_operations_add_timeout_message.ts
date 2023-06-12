import { MigrationInterface, QueryRunner } from "typeorm";

export class eventOperationsAddTimeoutMessage1686608540661 implements MigrationInterface {
    name = 'eventOperationsAddTimeoutMessage1686608540661'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ADD "timeoutMessage" text`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations" DROP COLUMN "timeoutMessage"`);
    }

}
