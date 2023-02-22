import { MigrationInterface, QueryRunner } from "typeorm";

export class commandsSetNullGroupOnDelete1677090770543 implements MigrationInterface {
    name = 'commandsSetNullGroupOnDelete1677090770543'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_03d418a239ea80b80ebc562999b"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_03d418a239ea80b80ebc562999b" FOREIGN KEY ("groupId") REFERENCES "channels_commands_groups"("id") ON DELETE SET NULL ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_03d418a239ea80b80ebc562999b"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_03d418a239ea80b80ebc562999b" FOREIGN KEY ("groupId") REFERENCES "channels_commands_groups"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

}
