import { MigrationInterface, QueryRunner } from "typeorm";

export class addObsWebsocketModule1676390859714 implements MigrationInterface {
    name = 'addObsWebsocketModule1676390859714'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TYPE "public"."channels_modules_settings_type_enum" RENAME TO "channels_modules_settings_type_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_modules_settings_type_enum" AS ENUM('youtube_song_requests', 'obs_websocket')`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "type" TYPE "public"."channels_modules_settings_type_enum" USING "type"::"text"::"public"."channels_modules_settings_type_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_modules_settings_type_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`CREATE TYPE "public"."channels_modules_settings_type_enum_old" AS ENUM('youtube_song_requests')`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "type" TYPE "public"."channels_modules_settings_type_enum_old" USING "type"::"text"::"public"."channels_modules_settings_type_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."channels_modules_settings_type_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_modules_settings_type_enum_old" RENAME TO "channels_modules_settings_type_enum"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
    }

}
