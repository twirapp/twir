import { MigrationInterface, QueryRunner } from "typeorm";

export class Giveaways1685365355180 implements MigrationInterface {
    name = 'Giveaways1685365355180'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."channels_giveaways_type_enum" AS ENUM('BY_KEYWORD', 'BY_RANDOM_NUMBER')`);
        await queryRunner.query(`CREATE TABLE "channels_giveaways" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "description" text NOT NULL, "type" "public"."channels_giveaways_type_enum" NOT NULL, "channel_id" text NOT NULL, "created_at" TIMESTAMP NOT NULL, "start_at" TIMESTAMP NOT NULL, "end_at" TIMESTAMP, "closed_at" TIMESTAMP NOT NULL, "required_min_watch_time" integer, "required_min_follow_time" integer, "required_min_messages" integer, "required_min_subscriber_tier" integer, "required_min_subscribe_time" integer, "eligible_user_groups" text NOT NULL, "keyword" character varying, "random_number_from" integer, "random_number_to" integer, "winning_random_number" integer, "winners_count" integer NOT NULL, "subscribers_luck" integer NOT NULL DEFAULT '0', "subscribers_tier1_luck" integer NOT NULL DEFAULT '0', "subscribers_tier2_luck" integer NOT NULL DEFAULT '0', "subscribers_tier3_luck" integer NOT NULL DEFAULT '0', "watched_time_lucks" text array NOT NULL DEFAULT '{}', "messages_lucks" text array NOT NULL DEFAULT '{}', "used_channel_points_lucks" text array NOT NULL DEFAULT '{}', CONSTRAINT "PK_cbdd6f63d58153d5ac295bdae16" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE TABLE "channels_giveaways_participants" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "giveaway_id" uuid NOT NULL, "is_winner" boolean NOT NULL DEFAULT false, "user_id" text NOT NULL, "is_subscriber" boolean NOT NULL DEFAULT false, "subscriber_tier" integer NOT NULL DEFAULT '1', "user_follow_since" TIMESTAMP, "user_stats_watched_time" bigint NOT NULL, "messages" integer NOT NULL DEFAULT '0', CONSTRAINT "PK_91852baf51c3862eb2215cceb12" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "FK_b5f1c883e497ba7a0eeae08e8b8"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "UQ_e74a3ef66bba62b18e3448211f7"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "UQ_b5f1c883e497ba7a0eeae08e8b8" UNIQUE ("userId")`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "UQ_e74a3ef66bba62b18e3448211f7" UNIQUE ("channelId", "userId")`);
        await queryRunner.query(`ALTER TABLE "channels_giveaways" ADD CONSTRAINT "FK_6b7386e0ef8c07770c2e12e627b" FOREIGN KEY ("channel_id") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "FK_eef18a143409b24bc98cd4b883d" FOREIGN KEY ("giveaway_id") REFERENCES "channels_giveaways"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "FK_9bff4bb1c2c8cfaee69bfba9fb8" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "FK_b5f1c883e497ba7a0eeae08e8b8" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "FK_b5f1c883e497ba7a0eeae08e8b8"`);
        await queryRunner.query(`ALTER TABLE "channels_giveaways_participants" DROP CONSTRAINT "FK_9bff4bb1c2c8cfaee69bfba9fb8"`);
        await queryRunner.query(`ALTER TABLE "channels_giveaways_participants" DROP CONSTRAINT "FK_eef18a143409b24bc98cd4b883d"`);
        await queryRunner.query(`ALTER TABLE "channels_giveaways" DROP CONSTRAINT "FK_6b7386e0ef8c07770c2e12e627b"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "UQ_e74a3ef66bba62b18e3448211f7"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "UQ_b5f1c883e497ba7a0eeae08e8b8"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "UQ_e74a3ef66bba62b18e3448211f7" UNIQUE ("channelId", "userId")`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "FK_b5f1c883e497ba7a0eeae08e8b8" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`DROP TABLE "channels_giveaways_participants"`);
        await queryRunner.query(`DROP TABLE "channels_giveaways"`);
        await queryRunner.query(`DROP TYPE "public"."channels_giveaways_type_enum"`);
    }

}
