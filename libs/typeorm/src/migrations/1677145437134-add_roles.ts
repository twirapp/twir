import { MigrationInterface, QueryRunner } from 'typeorm';

import { ChannelRoleType } from '../entities/ChannelRole';
import { RolePermissionEnum } from '../entities/RoleFlag';

export class addRoles1677145437134 implements MigrationInterface {
    name = 'addRoles1677145437134';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE "public"."channels_commands_permission_enum" RENAME TO "channels_commands_rolesids_enum"`);
        await queryRunner.query(`CREATE TYPE "public"."roles_permissions_permission_enum" AS ENUM('ADMINISTRATOR', 'UPDATE_CHANNEL_TITLE', 'UPDATE_CHANNEL_CATEGORY', 'VIEW_COMMANDS', 'MANAGE_COMMANDS', 'VIEW_KEYWORDS', 'MANAGE_KEYWORDS', 'VIEW_TIMERS', 'MANAGE_TIMERS', 'VIEW_INTEGRATIONS', 'MANAGE_INTEGRATIONS', 'VIEW_SONG_REQUESTS', 'MANAGE_SONG_REQUESTS', 'VIEW_MODERATION', 'MANAGE_MODERATION', 'VIEW_VARIABLES', 'MANAGE_VARIABLES', 'VIEW_GREETINGS', 'MANAGE_GREETINGS')`);
        await queryRunner.query(`CREATE TABLE "roles_permissions" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "permission" "public"."roles_permissions_permission_enum" NOT NULL DEFAULT 'ADMINISTRATOR', CONSTRAINT "UQ_2adf6a8f8f904e32f306f82579a" UNIQUE ("permission"), CONSTRAINT "PK_298f2c0e2ea45289aa0c4ac8a02" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE TABLE "channels_roles_permissions" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "roleId" uuid NOT NULL, "permissionId" uuid NOT NULL, CONSTRAINT "PK_8d68a0ba1ed1d48ed25833a2ec4" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE TYPE "public"."channels_roles_type_enum" AS ENUM('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'CUSTOM')`);
        await queryRunner.query(`CREATE TABLE "channels_roles" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "channelId" text NOT NULL, "name" character varying NOT NULL, "type" "public"."channels_roles_type_enum" NOT NULL DEFAULT 'CUSTOM', "system" boolean NOT NULL DEFAULT false, CONSTRAINT "PK_8445a7f7c362bbebcaab4de9845" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE TABLE "channel_role_users" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "userId" text NOT NULL, "roleId" uuid NOT NULL, CONSTRAINT "PK_fc2d9f175a21713ed41d1327e8e" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "rolesIds" text array NOT NULL DEFAULT '{}'`);

        for (const permission of Object.values(RolePermissionEnum)) {
            await queryRunner.query(`INSERT INTO "roles_permissions" ("permission") VALUES ($1)`, [permission]);
        }

        const administratorPermission = await queryRunner.query(`SELECT * FROM "roles_permissions" WHERE "permission" = $1`, [RolePermissionEnum.CAN_ACCESS_DASHBOARD]);

        for (const channel of await queryRunner.query(`SELECT * FROM "channels"`)) {
            for (const role of Object.values(ChannelRoleType).filter(v => v != ChannelRoleType.CUSTOM)) {
                await queryRunner.query('INSERT INTO "channels_roles" ("channelId", "name", "type", "system") VALUES ($1, $2, $3, $4)', [
                    channel.id,
                    role.charAt(0).toUpperCase() + role.slice(1).toLowerCase(),
                    role,
                    true,
                ]);
            }

            const broadcasterRole = await queryRunner.query(`SELECT * FROM "channels_roles" WHERE "channelId" = $1 AND "type" = $2`, [channel.id, ChannelRoleType.BROADCASTER]);
            await queryRunner.query(`INSERT INTO "channels_roles_permissions" ("roleId", "permissionId") VALUES ($1, $2)`, [broadcasterRole[0].id, administratorPermission[0].id]);

            const roles = await queryRunner.query(`SELECT * FROM "channels_roles" WHERE "channelId" = $1`, [channel.id]);
            const commands = await queryRunner.query(`SELECT * FROM "channels_commands"`);

            for (const command of commands) {
                const permission = command.permission;
                const role = roles.find((r: any) => r.type === permission);
                if (role) {
                    await queryRunner.query(`update "channels_commands" set "rolesIds" = array_append("rolesIds", $1) where "id" = $2`, [role.id, command.id]);
                }
            }
        }

        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "permission"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" ADD CONSTRAINT "FK_6aa26bd2af43cc1f541cb5fd1b3" FOREIGN KEY ("roleId") REFERENCES "channels_roles"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" ADD CONSTRAINT "FK_fbe7180b62337b4479a0939388a" FOREIGN KEY ("permissionId") REFERENCES "roles_permissions"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles" ADD CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channel_role_users" ADD CONSTRAINT "FK_1fc5209efb1fdf67c107c2ed4cc" FOREIGN KEY ("roleId") REFERENCES "channels_roles"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channel_role_users" ADD CONSTRAINT "FK_e8a94b3e9e4fa65a874a318ac3f" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channel_role_users" DROP CONSTRAINT "FK_e8a94b3e9e4fa65a874a318ac3f"`);
        await queryRunner.query(`ALTER TABLE "channel_role_users" DROP CONSTRAINT "FK_1fc5209efb1fdf67c107c2ed4cc"`);
        await queryRunner.query(`ALTER TABLE "channels_roles" DROP CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" DROP CONSTRAINT "FK_fbe7180b62337b4479a0939388a"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_permissions" DROP CONSTRAINT "FK_6aa26bd2af43cc1f541cb5fd1b3"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "rolesIds"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "rolesIds" "public"."channels_commands_rolesids_enum" NOT NULL`);
        await queryRunner.query(`DROP TABLE "channel_role_users"`);
        await queryRunner.query(`DROP TABLE "channels_roles"`);
        await queryRunner.query(`DROP TYPE "public"."channels_roles_type_enum"`);
        await queryRunner.query(`DROP TABLE "channels_roles_permissions"`);
        await queryRunner.query(`DROP TABLE "roles_permissions"`);
        await queryRunner.query(`DROP TYPE "public"."roles_permissions_permission_enum"`);
        await queryRunner.query(`ALTER TYPE "public"."channels_commands_rolesids_enum" RENAME TO "channels_commands_permission_enum"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" RENAME COLUMN "rolesIds" TO "permission"`);
    }

}
