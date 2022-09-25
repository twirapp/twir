import { MigrationInterface, QueryRunner } from 'typeorm';

export class addDonationEvent1664050030373 implements MigrationInterface {
  name = 'addDonationEvent1664050030373';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "channel_events_donations" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "eventId" uuid NOT NULL, "fromUserId" character varying, "toUserId" character varying, "amount" numeric NOT NULL, "currency" character varying NOT NULL, "username" character varying, "message" character varying, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), CONSTRAINT "REL_e682e4bbfe3a354505755b7ee9" UNIQUE ("eventId"), CONSTRAINT "PK_968741787daf46977ee9e0f23c8" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_donations" ADD CONSTRAINT "FK_e682e4bbfe3a354505755b7ee9b" FOREIGN KEY ("eventId") REFERENCES "channel_events_list"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channel_events_donations" DROP CONSTRAINT "FK_e682e4bbfe3a354505755b7ee9b"`,
    );
    await queryRunner.query(`DROP TABLE "channel_events_donations"`);
  }
}
