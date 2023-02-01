import { MigrationInterface, QueryRunner } from 'typeorm';

export class addMissedServiceToEnum1675262407192 implements MigrationInterface {
    name = 'addMissedServiceToEnum1675262407192';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum" RENAME TO "integrations_service_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY', 'DONATELLO')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`INSERT INTO public.integrations (id, service, "accessToken", "refreshToken", "clientId", "clientSecret", "apiKey", "redirectUrl") VALUES ('8d13c58a-a3a3-4b2d-8578-a6ac10ed2a48', 'DONATELLO', null, null, null, null, null, null)`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ALTER COLUMN "id" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum_old" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum_old" USING "service"::"text"::"public"."integrations_service_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum_old" RENAME TO "integrations_service_enum"`);
    }

}
