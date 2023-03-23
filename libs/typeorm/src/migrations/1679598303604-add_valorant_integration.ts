import { MigrationInterface, QueryRunner } from 'typeorm';

export class addValorantIntegration1679598303604 implements MigrationInterface {
    name = 'addValorantIntegration1679598303604';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum" RENAME TO "integrations_service_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY', 'DONATELLO', 'VALORANT')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum_old"`);
        await queryRunner.query(`INSERT INTO "integrations" ("service") VALUES ($1)`, ['VALORANT']);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum_old" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY', 'DONATELLO')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum_old" USING "service"::"text"::"public"."integrations_service_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum_old" RENAME TO "integrations_service_enum"`);
    }

}
