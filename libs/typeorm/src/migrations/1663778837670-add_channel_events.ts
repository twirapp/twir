import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelEvents1663778837670 implements MigrationInterface {
  name = 'addChannelEvents1663778837670';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TYPE "public"."channel_events_list_type_enum" AS ENUM('follow', 'subscription', 'resubscription', 'donation', 'host', 'raid', 'moderator_added', 'moderator_remove')`,
    );
    await queryRunner.query(
      `CREATE TABLE "channel_events_list" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "type" "public"."channel_events_list_type_enum" NOT NULL, "channelId" text, CONSTRAINT "PK_ecc7dc4d42ccab404c461e25a4d" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `CREATE TABLE "channel_events_follows" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "eventId" uuid NOT NULL, "fromUserId" character varying NOT NULL, "toUserId" character varying NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), CONSTRAINT "REL_4420e6ba1604749ff3184fefe0" UNIQUE ("eventId"), CONSTRAINT "PK_0712d150127df2f1bc9d3ebd747" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" ADD CONSTRAINT "FK_62386183f7a7575141594880b03" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_follows" ADD CONSTRAINT "FK_4420e6ba1604749ff3184fefe01" FOREIGN KEY ("eventId") REFERENCES "channel_events_list"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channel_events_follows" DROP CONSTRAINT "FK_4420e6ba1604749ff3184fefe01"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" DROP CONSTRAINT "FK_62386183f7a7575141594880b03"`,
    );
    await queryRunner.query(`DROP TABLE "channel_events_follows"`);
    await queryRunner.query(`DROP TABLE "channel_events_list"`);
    await queryRunner.query(`DROP TYPE "public"."channel_events_list_type_enum"`);
  }
}
