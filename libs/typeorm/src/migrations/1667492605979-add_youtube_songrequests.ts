import { MigrationInterface, QueryRunner } from 'typeorm';

export class addYoutubeSongrequests1667492605979 implements MigrationInterface {
  name = 'addYoutubeSongrequests1667492605979';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TABLE "channels_requested_songs" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "videoId" character varying NOT NULL, "title" text NOT NULL, "duration" integer NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "orderedById" text NOT NULL, "channelId" text NOT NULL, CONSTRAINT "PK_5f1d2b4311bd53a7cd499f1b59f" PRIMARY KEY ("id"))`,
    );
    await queryRunner.query(
      `ALTER TYPE "public"."channels_commands_module_enum" RENAME TO "channels_commands_module_enum_old"`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_module_enum" AS ENUM('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION', 'MANAGE', 'SONGREQUEST')`,
    );
    await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum" USING "module"::"text"::"public"."channels_commands_module_enum"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum_old"`);
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "blackListSentences" DROP NOT NULL`,
    );
    await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "users" ALTER COLUMN "apiKey" SET DEFAULT gen_random_uuid()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ADD CONSTRAINT "FK_bf976a08dee5c9ffa7e1773defe" FOREIGN KEY ("orderedById") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" ADD CONSTRAINT "FK_a757e3014566676c024e4ce16d1" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`,
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" DROP CONSTRAINT "FK_a757e3014566676c024e4ce16d1"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_requested_songs" DROP CONSTRAINT "FK_bf976a08dee5c9ffa7e1773defe"`,
    );
    await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "users" ALTER COLUMN "apiKey" SET DEFAULT uuid_generate_v4()`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_moderation_settings" ALTER COLUMN "blackListSentences" SET NOT NULL`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_module_enum_old" AS ENUM('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION', 'MANAGE')`,
    );
    await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum_old" USING "module"::"text"::"public"."channels_commands_module_enum_old"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum"`);
    await queryRunner.query(
      `ALTER TYPE "public"."channels_commands_module_enum_old" RENAME TO "channels_commands_module_enum"`,
    );
    await queryRunner.query(`DROP TABLE "channels_requested_songs"`);
  }
}
