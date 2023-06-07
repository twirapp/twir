import { MigrationInterface, QueryRunner } from "typeorm";

export class addRocketleagueIntegration1686119058914 implements MigrationInterface {
    name = 'addRocketleagueIntegration1686119058914'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum" RENAME TO "integrations_service_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS', 'DONATEPAY', 'DONATELLO', 'VALORANT', 'DONATE_STREAM', 'ROCKETLEAGUE')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum_old"`);
        await queryRunner.query(`INSERT INTO "integrations" ("service") VALUES ($1)`, ['ROCKETLEAGUE']);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."integrations_service_enum_old" AS ENUM('DONATELLO', 'DONATEPAY', 'DONATE_STREAM', 'DONATIONALERTS', 'FACEIT', 'LASTFM', 'SPOTIFY', 'STREAMLABS', 'VALORANT', 'VK')`);
        await queryRunner.query(`ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum_old" USING "service"::"text"::"public"."integrations_service_enum_old"`);
        await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."integrations_service_enum_old" RENAME TO "integrations_service_enum"`);
    }

}
