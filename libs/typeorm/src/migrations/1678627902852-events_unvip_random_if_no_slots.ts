import { MigrationInterface, QueryRunner } from 'typeorm';

export class eventsUnvipRandomIfNoSlots1678627902852 implements MigrationInterface {
    name = 'eventsUnvipRandomIfNoSlots1678627902852';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."channels_events_operations_type_enum" RENAME TO "channels_events_operations_type_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_events_operations_type_enum" AS ENUM('TIMEOUT', 'TIMEOUT_RANDOM', 'BAN', 'UNBAN', 'BAN_RANDOM', 'VIP', 'UNVIP', 'UNVIP_RANDOM', 'UNVIP_RANDOM_IF_NO_SLOTS', 'MOD', 'UNMOD', 'UNMOD_RANDOM', 'SEND_MESSAGE', 'CHANGE_TITLE', 'CHANGE_CATEGORY', 'FULFILL_REDEMPTION', 'CANCEL_REDEMPTION', 'ENABLE_SUBMODE', 'DISABLE_SUBMODE', 'ENABLE_EMOTEONLY', 'DISABLE_EMOTEONLY', 'CREATE_GREETING', 'OBS_SET_SCENE', 'OBS_TOGGLE_SOURCE', 'OBS_TOGGLE_AUDIO', 'OBS_AUDIO_SET_VOLUME', 'OBS_AUDIO_INCREASE_VOLUME', 'OBS_AUDIO_DECREASE_VOLUME', 'OBS_DISABLE_AUDIO', 'OBS_ENABLE_AUDIO', 'OBS_START_STREAM', 'OBS_STOP_STREAM', 'CHANGE_VARIABLE', 'INCREMENT_VARIABLE', 'DECREMENT_VARIABLE', 'TTS_SAY', 'TTS_SKIP', 'TTS_ENABLE', 'TTS_DISABLE', 'ALLOW_COMMAND_TO_USER', 'REMOVE_ALLOW_COMMAND_TO_USER', 'DENY_COMMAND_TO_USER', 'REMOVE_DENY_COMMAND_TO_USER')`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "type" TYPE "public"."channels_events_operations_type_enum" USING "type"::"text"::"public"."channels_events_operations_type_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_operations_type_enum_old"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."channels_events_operations_type_enum_old" AS ENUM('ALLOW_COMMAND_TO_USER', 'BAN', 'BAN_RANDOM', 'CANCEL_REDEMPTION', 'CHANGE_CATEGORY', 'CHANGE_TITLE', 'CHANGE_VARIABLE', 'CREATE_GREETING', 'DECREMENT_VARIABLE', 'DENY_COMMAND_TO_USER', 'DISABLE_EMOTEONLY', 'DISABLE_SUBMODE', 'ENABLE_EMOTEONLY', 'ENABLE_SUBMODE', 'FULFILL_REDEMPTION', 'INCREMENT_VARIABLE', 'MOD', 'OBS_AUDIO_DECREASE_VOLUME', 'OBS_AUDIO_INCREASE_VOLUME', 'OBS_AUDIO_SET_VOLUME', 'OBS_DISABLE_AUDIO', 'OBS_ENABLE_AUDIO', 'OBS_SET_SCENE', 'OBS_START_STREAM', 'OBS_STOP_STREAM', 'OBS_TOGGLE_AUDIO', 'OBS_TOGGLE_SOURCE', 'REMOVE_ALLOW_COMMAND_TO_USER', 'REMOVE_DENY_COMMAND_TO_USER', 'SEND_MESSAGE', 'TIMEOUT', 'TIMEOUT_RANDOM', 'TTS_DISABLE', 'TTS_ENABLE', 'TTS_SAY', 'TTS_SKIP', 'UNBAN', 'UNMOD', 'UNMOD_RANDOM', 'UNVIP', 'UNVIP_RANDOM', 'VIP')`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations" ALTER COLUMN "type" TYPE "public"."channels_events_operations_type_enum_old" USING "type"::"text"::"public"."channels_events_operations_type_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_operations_type_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_events_operations_type_enum_old" RENAME TO "channels_events_operations_type_enum"`);
    }

}
