import { MigrationInterface, QueryRunner } from "typeorm";

export class addCustomvarScriptType1677015345893 implements MigrationInterface {
    name = 'addCustomvarScriptType1677015345893'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."channels_customvars_type_enum" RENAME TO "channels_customvars_type_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."channels_customvars_type_enum" AS ENUM('SCRIPT', 'TEXT', 'NUMBER')`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "type" TYPE "public"."channels_customvars_type_enum" USING "type"::"text"::"public"."channels_customvars_type_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_customvars_type_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`CREATE TYPE "public"."channels_customvars_type_enum_old" AS ENUM('SCRIPT', 'TEXT')`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "type" TYPE "public"."channels_customvars_type_enum_old" USING "type"::"text"::"public"."channels_customvars_type_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."channels_customvars_type_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_customvars_type_enum_old" RENAME TO "channels_customvars_type_enum"`);
    }

}
