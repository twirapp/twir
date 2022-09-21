import { MigrationInterface, QueryRunner } from 'typeorm';

export class eventsAddCreatedAt1663779404997 implements MigrationInterface {
  name = 'eventsAddCreatedAt1663779404997';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" ADD "createdAt" TIMESTAMP NOT NULL DEFAULT now()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" DROP CONSTRAINT "FK_62386183f7a7575141594880b03"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" ALTER COLUMN "channelId" SET NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" ADD CONSTRAINT "FK_62386183f7a7575141594880b03" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" DROP CONSTRAINT "FK_62386183f7a7575141594880b03"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" ALTER COLUMN "channelId" DROP NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channel_events_list" ADD CONSTRAINT "FK_62386183f7a7575141594880b03" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
    await queryRunner.query(`ALTER TABLE "channel_events_list" DROP COLUMN "createdAt"`);
  }
}
