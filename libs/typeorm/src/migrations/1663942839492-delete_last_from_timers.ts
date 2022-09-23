import { MigrationInterface, QueryRunner } from "typeorm";

export class deleteLastFromTimers1663942839492 implements MigrationInterface {
    name = 'deleteLastFromTimers1663942839492'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_timers" DROP COLUMN "last"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_timers" ADD "last" integer NOT NULL DEFAULT '0'`);
    }

}
