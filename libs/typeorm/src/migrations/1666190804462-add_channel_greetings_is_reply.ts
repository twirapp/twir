import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelGreetingsIsReply1666190804462 implements MigrationInterface {
  name = 'addChannelGreetingsIsReply1666190804462';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_greetings" ADD "isReply" boolean NOT NULL DEFAULT true`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "channels_greetings" DROP COLUMN "isReply"`);
  }
}
