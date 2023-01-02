import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelModuleSettings1667580847530 implements MigrationInterface {
  name = 'addChannelModuleSettings1667580847530';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TYPE "public"."channels_modules_settings_type_enum" AS ENUM('youtube_song_requests')`,
    );
    await queryRunner.query(
      `CREATE TABLE "channels_modules_settings" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "type" "public"."channels_modules_settings_type_enum" NOT NULL, "settings" jsonb NOT NULL, "channelId" text NOT NULL, CONSTRAINT "PK_d5df4cbeb326c06be0e04654e36" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ADD "orderedByName" character varying NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "users" ALTER COLUMN "apiKey" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "FK_c145b2745bd936041f37b5d5d49" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "FK_c145b2745bd936041f37b5d5d49"`,
    );
    await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "users" ALTER COLUMN "apiKey" SET DEFAULT uuid_generate_v4()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`,
    );
    await queryRunner.query(`ALTER TABLE "channels_requested_songs" DROP COLUMN "orderedByName"`);
    await queryRunner.query(`DROP TABLE "channels_modules_settings"`);
    await queryRunner.query(`DROP TYPE "public"."channels_modules_settings_type_enum"`);
  }
}
