import { MigrationInterface, QueryRunner } from "typeorm";

export class addEvents1676039785645 implements MigrationInterface {
    name = 'addEvents1676039785645'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."channels_events_operations_type_enum" AS ENUM('BAN', 'UNBAN', 'BAN_RANDOM', 'VIP', 'UNVIP', 'UNVIP_RANDOM', 'MOD', 'UNMOD', 'UNMOD_RANDOM', 'SEND_MESSAGE', 'CHANGE_TITLE', 'CHANGE_CATEGORY', 'FULFILL_REDEMPTION', 'CANCEL_REDEMPTION', 'ENABLE_SUBMODE', 'DISABLE_SUBMODE', 'ENABLE_EMOTEONLY', 'DISABLE_EMOTEONLY')`);
        await queryRunner.query(`CREATE TABLE "channels_events_operations" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "type" "public"."channels_events_operations_type_enum" NOT NULL, "delay" integer, "eventId" uuid NOT NULL, "input" text, "repeat" integer, "order" integer NOT NULL, CONSTRAINT "PK_52cef67c53dd212f1e1b3621a61" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE TYPE "public"."channels_events_type_enum" AS ENUM('FOLLOW', 'SUBSCRIBE', 'RESUBSCRIBE', 'SUB_GIFT', 'REDEMPTION_CREATED', 'COMMAND_USED', 'FIRST_USER_MESSAGE', 'RAIDED', 'TITLE_OR_CATEGORY_CHANGED', 'STREAM_ONLINE', 'STREAM_OFFLINE', 'ON_CHAT_CLEAR', 'DONATE')`);
        await queryRunner.query(`CREATE TABLE "channels_events" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "type" "public"."channels_events_type_enum" NOT NULL, "description" text, "rewardId" uuid, "commandId" text, "channelId" text NOT NULL, CONSTRAINT "PK_9e9ce619150ad05221318513d4d" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ADD CONSTRAINT "FK_b2e27e84fa5bfbf8fd27ac9e948" FOREIGN KEY ("eventId") REFERENCES "channels_events"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_events" ADD CONSTRAINT "FK_763ec88e86ecbf8ca6ec3a9ec7b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events" DROP CONSTRAINT "FK_763ec88e86ecbf8ca6ec3a9ec7b"`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" DROP CONSTRAINT "FK_b2e27e84fa5bfbf8fd27ac9e948"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`DROP TABLE "channels_events"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_type_enum"`);
        await queryRunner.query(`DROP TABLE "channels_events_operations"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_operations_type_enum"`);
    }

}
