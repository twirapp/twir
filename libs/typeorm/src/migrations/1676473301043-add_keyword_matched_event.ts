import { MigrationInterface, QueryRunner } from "typeorm";

export class addKeywordMatchedEvent1676473301043 implements MigrationInterface {
    name = 'addKeywordMatchedEvent1676473301043'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."channels_events_type_enum" RENAME TO "channels_events_type_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_events_type_enum" AS ENUM('FOLLOW', 'SUBSCRIBE', 'RESUBSCRIBE', 'SUB_GIFT', 'REDEMPTION_CREATED', 'COMMAND_USED', 'FIRST_USER_MESSAGE', 'RAIDED', 'TITLE_OR_CATEGORY_CHANGED', 'STREAM_ONLINE', 'STREAM_OFFLINE', 'ON_CHAT_CLEAR', 'DONATE', 'KEYWORD_MATCHED')`);
        await queryRunner.query(`ALTER TABLE "channels_events" ALTER COLUMN "type" TYPE "public"."channels_events_type_enum" USING "type"::"text"::"public"."channels_events_type_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_type_enum_old"`);
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
        await queryRunner.query(`CREATE TYPE "public"."channels_events_type_enum_old" AS ENUM('FOLLOW', 'SUBSCRIBE', 'RESUBSCRIBE', 'SUB_GIFT', 'REDEMPTION_CREATED', 'COMMAND_USED', 'FIRST_USER_MESSAGE', 'RAIDED', 'TITLE_OR_CATEGORY_CHANGED', 'STREAM_ONLINE', 'STREAM_OFFLINE', 'ON_CHAT_CLEAR', 'DONATE')`);
        await queryRunner.query(`ALTER TABLE "channels_events" ALTER COLUMN "type" TYPE "public"."channels_events_type_enum_old" USING "type"::"text"::"public"."channels_events_type_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_type_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_events_type_enum_old" RENAME TO "channels_events_type_enum"`);
    }

}
