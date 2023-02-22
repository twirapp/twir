import { MigrationInterface, QueryRunner } from "typeorm";

export class commandGroupsAddColor1677101192628 implements MigrationInterface {
    name = 'commandGroupsAddColor1677101192628'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" ADD "color" character varying NOT NULL DEFAULT 'rgba(37, 38, 43, 1)'`);
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
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" DROP COLUMN "color"`);
    }

}
