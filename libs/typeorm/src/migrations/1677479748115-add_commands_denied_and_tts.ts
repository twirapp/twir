import { MigrationInterface, QueryRunner } from "typeorm";

export class addCommandsDeniedAndTts1677479748115 implements MigrationInterface {
    name = 'addCommandsDeniedAndTts1677479748115'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD "userId" text`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "UQ_b5f1c883e497ba7a0eeae08e8b8" UNIQUE ("userId")`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "deniedUsersIds" text array NOT NULL DEFAULT '{}'`);
        await queryRunner.query(`ALTER TYPE "public"."channels_modules_settings_type_enum" RENAME TO "channels_modules_settings_type_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_modules_settings_type_enum" AS ENUM('youtube_song_requests', 'obs_websocket', 'tts')`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "type" TYPE "public"."channels_modules_settings_type_enum" USING "type"::"text"::"public"."channels_modules_settings_type_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_modules_settings_type_enum_old"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_commands_module_enum" RENAME TO "channels_commands_module_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_commands_module_enum" AS ENUM('CUSTOM', 'DOTA', 'MODERATION', 'MANAGE', 'SONGREQUEST', 'TTS')`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum" USING "module"::"text"::"public"."channels_commands_module_enum"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`);
        await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "FK_b5f1c883e497ba7a0eeae08e8b8" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "FK_b5f1c883e497ba7a0eeae08e8b8"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`CREATE TYPE "public"."channels_commands_module_enum_old" AS ENUM('CUSTOM', 'DOTA', 'MANAGE', 'MODERATION', 'SONGREQUEST')`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum_old" USING "module"::"text"::"public"."channels_commands_module_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`);
        await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_commands_module_enum_old" RENAME TO "channels_commands_module_enum"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_modules_settings_type_enum_old" AS ENUM('obs_websocket', 'youtube_song_requests')`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "type" TYPE "public"."channels_modules_settings_type_enum_old" USING "type"::"text"::"public"."channels_modules_settings_type_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."channels_modules_settings_type_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_modules_settings_type_enum_old" RENAME TO "channels_modules_settings_type_enum"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "deniedUsersIds"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "UQ_b5f1c883e497ba7a0eeae08e8b8"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP COLUMN "userId"`);
    }

}
