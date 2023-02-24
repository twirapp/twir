import { MigrationInterface, QueryRunner } from 'typeorm';

export class renameChannelRoleUser21677196001685 implements MigrationInterface {
    name = 'renameChannelRoleUser21677196001685';

    public async up(queryRunner: QueryRunner): Promise<void> {
       // rename
        await queryRunner.query(`ALTER TABLE "channels_role_users" RENAME TO "channels_roles_users"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {

    }

}
