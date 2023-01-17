import { MigrationInterface, QueryRunner } from "typeorm";

export class addDonatePay1673969905358 implements MigrationInterface {
    name = 'addDonatePay1673969905358'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum" RENAME TO "integrations_service_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "users_online" ADD CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "users_online" ADD CONSTRAINT "FK_e40473bd90abb17377f9dedb12a" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "users_online" DROP CONSTRAINT "FK_e40473bd90abb17377f9dedb12a"`);
        await queryRunner.query(`ALTER TABLE "users_online" DROP CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum_old" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum_old" USING "service"::"text"::"public"."integrations_service_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum_old" RENAME TO "integrations_service_enum"`);
    }

}
