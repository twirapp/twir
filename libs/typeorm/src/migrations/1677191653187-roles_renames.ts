import { MigrationInterface, QueryRunner } from 'typeorm';

export class rolesRenames1677191653187 implements MigrationInterface {
    name = 'rolesRenames1677191653187';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" DROP CONSTRAINT "FK_fbe7180b62337b4479a0939388a"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" RENAME COLUMN "permissionId" TO "flagId"`);
        // rename table
        await queryRunner.query(`ALTER TABLE "roles_permissions" RENAME TO "roles_flags"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" ADD CONSTRAINT "FK_9c973968dbc7625c55b6cd540af" FOREIGN KEY ("flagId") REFERENCES "roles_flags"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" DROP CONSTRAINT "FK_9c973968dbc7625c55b6cd540af"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" RENAME COLUMN "flagId" TO "permissionId"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" ADD CONSTRAINT "FK_fbe7180b62337b4479a0939388a" FOREIGN KEY ("permissionId") REFERENCES "roles_permissions"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

}
