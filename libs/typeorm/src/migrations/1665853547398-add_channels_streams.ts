import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelsStreams1665853547398 implements MigrationInterface {
  name = 'addChannelsStreams1665853547398';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "channels_streams" ("id" character varying NOT NULL, "userId" text NOT NULL, "userLogin" text NOT NULL, "userName" text NOT NULL, "gameId" integer NOT NULL, "gameName" text NOT NULL, "communityIds" text array DEFAULT '{}', "type" text NOT NULL, "title" text NOT NULL, "viewerCount" integer NOT NULL, "startedAt" TIMESTAMP NOT NULL, "language" text NOT NULL, "thumbnailUrl" text NOT NULL, "tagIds" text array DEFAULT '{}', "isMature" boolean NOT NULL, "parsedMessages" integer NOT NULL DEFAULT '0', CONSTRAINT "PK_d2b19f073bd39c68e9a9c6cca32" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_streams" ADD CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa" FOREIGN KEY ("userId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_streams" DROP CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa"`,
    );
    await queryRunner.query(`DROP TABLE "channels_streams"`);
  }
}
