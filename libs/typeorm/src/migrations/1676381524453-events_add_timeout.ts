import { MigrationInterface, QueryRunner } from "typeorm";

export class eventsAddTimeout1676381524453 implements MigrationInterface {
    name = 'eventsAddTimeout1676381524453'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ADD "timeoutTime" integer NOT NULL DEFAULT '600'`);
        await queryRunner.query(`ALTER TYPE "public"."channels_events_operations_type_enum" RENAME TO "channels_events_operations_type_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_events_operations_type_enum" AS ENUM('TIMEOUT', 'BAN', 'UNBAN', 'BAN_RANDOM', 'VIP', 'UNVIP', 'UNVIP_RANDOM', 'MOD', 'UNMOD', 'UNMOD_RANDOM', 'SEND_MESSAGE', 'CHANGE_TITLE', 'CHANGE_CATEGORY', 'FULFILL_REDEMPTION', 'CANCEL_REDEMPTION', 'ENABLE_SUBMODE', 'DISABLE_SUBMODE', 'ENABLE_EMOTEONLY', 'DISABLE_EMOTEONLY', 'CREATE_GREETING')`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "type" TYPE "public"."channels_events_operations_type_enum" USING "type"::"text"::"public"."channels_events_operations_type_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_operations_type_enum_old"`);
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
        await queryRunner.query(`CREATE TYPE "public"."channels_events_operations_type_enum_old" AS ENUM('BAN', 'UNBAN', 'BAN_RANDOM', 'VIP', 'UNVIP', 'UNVIP_RANDOM', 'MOD', 'UNMOD', 'UNMOD_RANDOM', 'SEND_MESSAGE', 'CHANGE_TITLE', 'CHANGE_CATEGORY', 'FULFILL_REDEMPTION', 'CANCEL_REDEMPTION', 'ENABLE_SUBMODE', 'DISABLE_SUBMODE', 'ENABLE_EMOTEONLY', 'DISABLE_EMOTEONLY', 'CREATE_GREETING')`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "type" TYPE "public"."channels_events_operations_type_enum_old" USING "type"::"text"::"public"."channels_events_operations_type_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_operations_type_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_events_operations_type_enum_old" RENAME TO "channels_events_operations_type_enum"`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" DROP COLUMN "timeoutTime"`);
    }

}
