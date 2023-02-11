import { MigrationInterface, QueryRunner } from "typeorm";

export class eventsMakeDelayZeroByDefault1676121683251 implements MigrationInterface {
    name = 'eventsMakeDelayZeroByDefault1676121683251'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "delay" SET NOT NULL`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "delay" SET DEFAULT '0'`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "delay" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "delay" DROP NOT NULL`);
    }

}
