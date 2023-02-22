import { MigrationInterface, QueryRunner } from "typeorm";

export class addCommandsGroups1677086846320 implements MigrationInterface {
    name = 'addCommandsGroups1677086846320'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "channels_commands_groups" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "channelId" text NOT NULL, "name" character varying NOT NULL, CONSTRAINT "PK_0e3eb1b93ec7980d98a87d4749a" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "groupId" uuid`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" ADD CONSTRAINT "FK_c202e2ed66394a1bd6651734078" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_03d418a239ea80b80ebc562999b" FOREIGN KEY ("groupId") REFERENCES "channels_commands_groups"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_03d418a239ea80b80ebc562999b"`);
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" DROP CONSTRAINT "FK_c202e2ed66394a1bd6651734078"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "groupId"`);
        await queryRunner.query(`DROP TABLE "channels_commands_groups"`);
    }

}
