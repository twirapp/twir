import { MigrationInterface, QueryRunner } from "typeorm";

export class rolesPermissionsRenamePermissionColumnToFlag1677188520301 implements MigrationInterface {
    name = 'rolesPermissionsRenamePermissionColumnToFlag1677188520301'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "roles_permissions" RENAME COLUMN "permission" TO "flag"`);
        await queryRunner.query(`ALTER TYPE "public"."roles_permissions_permission_enum" RENAME TO "roles_permissions_flag_enum"`);
        await queryRunner.query(`ALTER TABLE "roles_permissions" RENAME CONSTRAINT "UQ_2adf6a8f8f904e32f306f82579a" TO "UQ_aa1091e35aab7d9ae18694d84c0"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "roles_permissions" RENAME CONSTRAINT "UQ_aa1091e35aab7d9ae18694d84c0" TO "UQ_2adf6a8f8f904e32f306f82579a"`);
        await queryRunner.query(`ALTER TYPE "public"."roles_permissions_flag_enum" RENAME TO "roles_permissions_permission_enum"`);
        await queryRunner.query(`ALTER TABLE "roles_permissions" RENAME COLUMN "flag" TO "permission"`);
    }

}
