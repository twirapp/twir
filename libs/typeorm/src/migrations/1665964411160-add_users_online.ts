import { MigrationInterface, QueryRunner } from 'typeorm';

export class addUsersOnline1665964411160 implements MigrationInterface {
  name = 'addUsersOnline1665964411160';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "users_online" ("id" text NOT NULL DEFAULT gen_random_uuid(), "channelId" text NOT NULL, "userId" text, "userName" text, CONSTRAINT "PK_efb961b9e7ba28a837498e14827" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_online" ADD CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_online" ADD CONSTRAINT "FK_e40473bd90abb17377f9dedb12a" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "users_online" DROP CONSTRAINT "FK_e40473bd90abb17377f9dedb12a"`,
    );
    await queryRunner.query(
      `ALTER TABLE "users_online" DROP CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4"`,
    );
    await queryRunner.query(`DROP TABLE "users_online"`);
  }
}
