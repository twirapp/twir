import { MigrationInterface, QueryRunner } from 'typeorm';

export class changeIsReplyForSpam1664810074334 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`UPDATE "channels_commands" SET "is_reply"=$1 WHERE "defaultName"=$2`, [
      false,
      'spam',
    ]);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    return;
  }
}
