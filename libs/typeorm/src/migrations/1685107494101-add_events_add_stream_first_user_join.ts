import { MigrationInterface, QueryRunner } from 'typeorm';

export class addEventsAddStreamFirstUserJoin1685107494101 implements MigrationInterface {
	name = 'addEventsAddStreamFirstUserJoin1685107494101';

	public async up(queryRunner: QueryRunner): Promise<void> {
		await queryRunner.query(
			`ALTER TYPE "public"."channels_events_type_enum" RENAME TO "channels_events_type_enum_old"`,
		);
		await queryRunner.query(
			`CREATE TYPE "public"."channels_events_type_enum" AS ENUM('FOLLOW', 'SUBSCRIBE', 'RESUBSCRIBE', 'SUB_GIFT', 'REDEMPTION_CREATED', 'COMMAND_USED', 'FIRST_USER_MESSAGE', 'RAIDED', 'TITLE_OR_CATEGORY_CHANGED', 'STREAM_ONLINE', 'STREAM_OFFLINE', 'ON_CHAT_CLEAR', 'DONATE', 'KEYWORD_MATCHED', 'GREETING_SENDED', 'POLL_BEGIN', 'POLL_PROGRESS', 'POLL_END', 'PREDICTION_BEGIN', 'PREDICTION_PROGRESS', 'PREDICTION_END', 'PREDICTION_LOCK', 'STREAM_FIRST_USER_JOIN')`,
		);
		await queryRunner.query(
			`ALTER TABLE "channels_events" ALTER COLUMN "type" TYPE "public"."channels_events_type_enum" USING "type"::"text"::"public"."channels_events_type_enum"`,
		);
		await queryRunner.query(`DROP TYPE "public"."channels_events_type_enum_old"`);
	}

	public async down(queryRunner: QueryRunner): Promise<void> {
		await queryRunner.query(
			`CREATE TYPE "public"."channels_events_type_enum_old" AS ENUM('FOLLOW', 'SUBSCRIBE', 'RESUBSCRIBE', 'SUB_GIFT', 'REDEMPTION_CREATED', 'COMMAND_USED', 'FIRST_USER_MESSAGE', 'RAIDED', 'TITLE_OR_CATEGORY_CHANGED', 'STREAM_ONLINE', 'STREAM_OFFLINE', 'ON_CHAT_CLEAR', 'DONATE', 'KEYWORD_MATCHED', 'GREETING_SENDED', 'POLL_BEGIN', 'POLL_PROGRESS', 'POLL_END', 'PREDICTION_BEGIN', 'PREDICTION_PROGRESS', 'PREDICTION_END', 'PREDICTION_LOCK')`,
		);
		await queryRunner.query(
			`ALTER TABLE "channels_events" ALTER COLUMN "type" TYPE "public"."channels_events_type_enum_old" USING "type"::"text"::"public"."channels_events_type_enum_old"`,
		);
		await queryRunner.query(`DROP TYPE "public"."channels_events_type_enum"`);
		await queryRunner.query(
			`ALTER TYPE "public"."channels_events_type_enum_old" RENAME TO "channels_events_type_enum"`,
		);
	}
}
