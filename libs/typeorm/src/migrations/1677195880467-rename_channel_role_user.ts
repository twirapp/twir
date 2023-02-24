import { MigrationInterface, QueryRunner } from 'typeorm';

export class renameChannelRoleUser1677195880467 implements MigrationInterface {
    name = 'renameChannelRoleUser1677195880467';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "channels_role_users" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "userId" text NOT NULL, "roleId" uuid NOT NULL, CONSTRAINT "PK_1dcc4d6efe4a9c845a435a13157" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TYPE "public"."roles_permissions_flag_enum" RENAME TO "roles_permissions_flag_enum_old"`);
        await queryRunner.query(`CREATE TYPE "public"."roles_flags_flag_enum" AS ENUM('ADMINISTRATOR', 'UPDATE_CHANNEL_TITLE', 'UPDATE_CHANNEL_CATEGORY', 'VIEW_COMMANDS', 'MANAGE_COMMANDS', 'VIEW_KEYWORDS', 'MANAGE_KEYWORDS', 'VIEW_TIMERS', 'MANAGE_TIMERS', 'VIEW_INTEGRATIONS', 'MANAGE_INTEGRATIONS', 'VIEW_SONG_REQUESTS', 'MANAGE_SONG_REQUESTS', 'VIEW_MODERATION', 'MANAGE_MODERATION', 'VIEW_VARIABLES', 'MANAGE_VARIABLES', 'VIEW_GREETINGS', 'MANAGE_GREETINGS')`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" TYPE "public"."roles_flags_flag_enum" USING "flag"::"text"::"public"."roles_flags_flag_enum"`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" SET DEFAULT 'ADMINISTRATOR'`);
        await queryRunner.query(`DROP TYPE "public"."roles_permissions_flag_enum_old"`);
        await queryRunner.query(`ALTER TABLE "channels_role_users" ADD CONSTRAINT "FK_111e0f052f520d82c37d5348dfc" FOREIGN KEY ("roleId") REFERENCES "channels_roles"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_role_users" ADD CONSTRAINT "FK_89918fada6711aa92e2e6401d63" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`DROP TABLE "channel_role_users"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_role_users" DROP CONSTRAINT "FK_89918fada6711aa92e2e6401d63"`);
        await queryRunner.query(`ALTER TABLE "channels_role_users" DROP CONSTRAINT "FK_111e0f052f520d82c37d5348dfc"`);
        await queryRunner.query(`CREATE TYPE "public"."roles_permissions_flag_enum_old" AS ENUM('ADMINISTRATOR', 'MANAGE_COMMANDS', 'MANAGE_GREETINGS', 'MANAGE_INTEGRATIONS', 'MANAGE_KEYWORDS', 'MANAGE_MODERATION', 'MANAGE_SONG_REQUESTS', 'MANAGE_TIMERS', 'MANAGE_VARIABLES', 'UPDATE_CHANNEL_CATEGORY', 'UPDATE_CHANNEL_TITLE', 'VIEW_COMMANDS', 'VIEW_GREETINGS', 'VIEW_INTEGRATIONS', 'VIEW_KEYWORDS', 'VIEW_MODERATION', 'VIEW_SONG_REQUESTS', 'VIEW_TIMERS', 'VIEW_VARIABLES')`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" TYPE "public"."roles_permissions_flag_enum_old" USING "flag"::"text"::"public"."roles_permissions_flag_enum_old"`);
        await queryRunner.query(`ALTER TABLE "roles_flags" ALTER COLUMN "flag" SET DEFAULT 'ADMINISTRATOR'`);
        await queryRunner.query(`DROP TYPE "public"."roles_flags_flag_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."roles_permissions_flag_enum_old" RENAME TO "roles_permissions_flag_enum"`);
        await queryRunner.query(`DROP TABLE "channels_role_users"`);
    }

}
