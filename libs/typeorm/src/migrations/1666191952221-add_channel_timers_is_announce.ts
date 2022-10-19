import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelTimersIsAnnounce1666191952221 implements MigrationInterface {
  name = 'addChannelTimersIsAnnounce1666191952221';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "channels_timers_responses" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "text" text NOT NULL, "isAnnounce" boolean NOT NULL DEFAULT true, "timerId" uuid NOT NULL, CONSTRAINT "PK_9ab6dbe62d496cde6d8e52f11f9" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers_responses" ADD CONSTRAINT "FK_464cac4218d2ba9ac84325d06c8" FOREIGN KEY ("timerId") REFERENCES "channels_timers"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );

    const timers = await queryRunner.query(`SELECT * from "channels_timers"`);

    for (const timer of timers) {
      for (const response of timer.responses) {
        await queryRunner.query(
          `INSERT INTO "channels_timers_responses" ("text", "isAnnounce", "timerId")
            VALUES ($1, $2, $3)`,
          [response, true, timer.id],
        );
      }
    }

    await queryRunner.query(`ALTER TABLE "channels_timers" DROP COLUMN "responses"`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_timers_responses" DROP CONSTRAINT "FK_464cac4218d2ba9ac84325d06c8"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_timers" ADD "responses" text array NOT NULL DEFAULT '{}'`,
    );
    await queryRunner.query(`DROP TABLE "channels_timers_responses"`);
  }
}
