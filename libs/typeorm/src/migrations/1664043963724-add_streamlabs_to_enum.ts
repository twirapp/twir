import { MigrationInterface, QueryRunner } from 'typeorm';

export class addStreamlabsToEnum1664043963724 implements MigrationInterface {
  name = 'addStreamlabsToEnum1664043963724';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TYPE "public"."integrations_service_enum" RENAME TO "integrations_service_enum_old"`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."integrations_service_enum" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS', 'STREAMLABS')`,
    );
    await queryRunner.query(
      `ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum" USING "service"::"text"::"public"."integrations_service_enum"`,
    );
    await queryRunner.query(`DROP TYPE "public"."integrations_service_enum_old"`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TYPE "public"."integrations_service_enum_old" AS ENUM('LASTFM', 'VK', 'FACEIT', 'SPOTIFY', 'DONATIONALERTS')`,
    );
    await queryRunner.query(
      `ALTER TABLE "integrations" ALTER COLUMN "service" TYPE "public"."integrations_service_enum_old" USING "service"::"text"::"public"."integrations_service_enum_old"`,
    );
    await queryRunner.query(`DROP TYPE "public"."integrations_service_enum"`);
    await queryRunner.query(
      `ALTER TYPE "public"."integrations_service_enum_old" RENAME TO "integrations_service_enum"`,
    );
  }
}
