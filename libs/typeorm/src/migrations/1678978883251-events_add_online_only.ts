import { MigrationInterface, QueryRunner } from 'typeorm';

export class eventsAddOnlineOnly1678978883251 implements MigrationInterface {
    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events" ADD "online_only" boolean NOT NULL DEFAULT false`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events" DROP COLUMN "online_only"`);
    }
}
