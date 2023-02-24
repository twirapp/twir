import { MigrationInterface, QueryRunner } from "typeorm";

export class rolesRenameAdministrator1677209480948 implements MigrationInterface {
    name = 'rolesRenameAdministrator1677209480948'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles_users" DROP CONSTRAINT "FK_89918fada6711aa92e2e6401d63"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" DROP CONSTRAINT "FK_111e0f052f520d82c37d5348dfc"`);
        await queryRunner.query(`ALTER TYPE "public"."roles_flags_flag_enum" RENAME TO "roles_flags_flag_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."roles_flags_flag_enum" AS ENUM('CAN_ACCESS_DASHBOARD', 'UPDATE_CHANNEL_TITLE', 'UPDATE_CHANNEL_CATEGORY', 'VIEW_COMMANDS', 'MANAGE_COMMANDS', 'VIEW_KEYWORDS', 'MANAGE_KEYWORDS', 'VIEW_TIMERS', 'MANAGE_TIMERS', 'VIEW_INTEGRATIONS', 'MANAGE_INTEGRATIONS', 'VIEW_SONG_REQUESTS', 'MANAGE_SONG_REQUESTS', 'VIEW_MODERATION', 'MANAGE_MODERATION', 'VIEW_VARIABLES', 'MANAGE_VARIABLES', 'VIEW_GREETINGS', 'MANAGE_GREETINGS')`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" TYPE "public"."roles_flags_flag_enum" USING "flag"::"text"::"public"."roles_flags_flag_enum"`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" SET DEFAULT 'CAN_ACCESS_DASHBOARD'`);
        await queryRunner.query(`DROP TYPE "public"."roles_flags_flag_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" ADD CONSTRAINT "FK_ae2e54203d7fed7ed38243e8b13" FOREIGN KEY ("roleId") REFERENCES "channels_roles"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" ADD CONSTRAINT "FK_c3b3d3917f24ec3fe3f5049e667" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles_users" DROP CONSTRAINT "FK_c3b3d3917f24ec3fe3f5049e667"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" DROP CONSTRAINT "FK_ae2e54203d7fed7ed38243e8b13"`);
        await queryRunner.query(`CREATE TYPE "public"."roles_flags_flag_enum_old" AS ENUM('ADMINISTRATOR', 'UPDATE_CHANNEL_TITLE', 'UPDATE_CHANNEL_CATEGORY', 'VIEW_COMMANDS', 'MANAGE_COMMANDS', 'VIEW_KEYWORDS', 'MANAGE_KEYWORDS', 'VIEW_TIMERS', 'MANAGE_TIMERS', 'VIEW_INTEGRATIONS', 'MANAGE_INTEGRATIONS', 'VIEW_SONG_REQUESTS', 'MANAGE_SONG_REQUESTS', 'VIEW_MODERATION', 'MANAGE_MODERATION', 'VIEW_VARIABLES', 'MANAGE_VARIABLES', 'VIEW_GREETINGS', 'MANAGE_GREETINGS')`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" TYPE "public"."roles_flags_flag_enum_old" USING "flag"::"text"::"public"."roles_flags_flag_enum_old"`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" SET DEFAULT 'ADMINISTRATOR'`);
        await queryRunner.query(`DROP TYPE "public"."roles_flags_flag_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."roles_flags_flag_enum_old" RENAME TO "roles_flags_flag_enum"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" ADD CONSTRAINT "FK_111e0f052f520d82c37d5348dfc" FOREIGN KEY ("roleId") REFERENCES "channels_roles"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" ADD CONSTRAINT "FK_89918fada6711aa92e2e6401d63" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

}
