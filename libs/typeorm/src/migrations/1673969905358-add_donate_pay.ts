import { MigrationInterface, QueryRunner } from 'typeorm';

export class addDonatePay1673969905358 implements MigrationInterface {
    name = 'addDonatePay1673969905358';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum" RENAME TO "integrations_service_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum_old"`);
        await queryRunner.query(`INSERT INTO public.integrations (id, service, "accessToken", "refreshToken", "clientId", "clientSecret", "apiKey", "redirectUrl") VALUES ('c71e1fdd-8e24-4e98-b515-fcc2fa9abf73', 'DONATEPAY', null, null, null, null, null, null)`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "users_online" DROP CONSTRAINT "FK_e40473bd90abb17377f9dedb12a"`);
        await queryRunner.query(`ALTER TABLE "users_online" DROP CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4"`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum_old" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum_old" USING "service"::"text"::"public"."integrations_service_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum_old" RENAME TO "integrations_service_enum"`);
    }

}
