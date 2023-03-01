import { MigrationInterface, QueryRunner } from 'typeorm';

import { RoleFlags, RoleType } from '../entities/ChannelRole';

export class addRoles1677211345172 implements MigrationInterface {
    name = 'addRoles1677211345172';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."channels_roles_type_enum" AS ENUM('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'CUSTOM')`);
        await queryRunner.query(`CREATE TYPE "public"."channels_roles_permissions_enum" AS ENUM('CAN_ACCESS_DASHBOARD', 'UPDATE_CHANNEL_TITLE', 'UPDATE_CHANNEL_CATEGORY', 'VIEW_COMMANDS', 'MANAGE_COMMANDS', 'VIEW_KEYWORDS', 'MANAGE_KEYWORDS', 'VIEW_TIMERS', 'MANAGE_TIMERS', 'VIEW_INTEGRATIONS', 'MANAGE_INTEGRATIONS', 'VIEW_SONG_REQUESTS', 'MANAGE_SONG_REQUESTS', 'VIEW_MODERATION', 'MANAGE_MODERATION', 'VIEW_VARIABLES', 'MANAGE_VARIABLES', 'VIEW_GREETINGS', 'MANAGE_GREETINGS')`);
        await queryRunner.query(`CREATE TABLE "channels_roles" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "channelId" text NOT NULL, "name" character varying NOT NULL, "type" "public"."channels_roles_type_enum" NOT NULL DEFAULT 'CUSTOM', "system" boolean NOT NULL DEFAULT false, "permissions" "public"."channels_roles_permissions_enum" array NOT NULL DEFAULT '{}', CONSTRAINT "PK_8445a7f7c362bbebcaab4de9845" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE TABLE "channels_roles_users" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "userId" text NOT NULL, "roleId" uuid NOT NULL, CONSTRAINT "PK_22c498bd2efb082200bb30ad44e" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_roles" ADD CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" ADD CONSTRAINT "FK_ae2e54203d7fed7ed38243e8b13" FOREIGN KEY ("roleId") REFERENCES "channels_roles"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" ADD CONSTRAINT "FK_c3b3d3917f24ec3fe3f5049e667" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);

        for (const channel of await queryRunner.query(`SELECT * FROM "channels"`)) {
            for (const role of Object.values(RoleType).filter(v => v != RoleType.CUSTOM)) {
                await queryRunner.query('INSERT INTO "channels_roles" ("channelId", "name", "type", "system", "permissions") VALUES ($1, $2, $3, $4, $5)', [
                    channel.id,
                    role.charAt(0).toUpperCase() + role.slice(1).toLowerCase(),
                    role,
                    true,
                    role === RoleType.BROADCASTER ? [RoleFlags.CAN_ACCESS_DASHBOARD] : [],
                ]);
            }

            const roles = await queryRunner.query(`SELECT * FROM "channels_roles" WHERE "channelId" = $1`, [channel.id]);
            const commands = await queryRunner.query(`SELECT * FROM "channels_commands" WHERE "channelId" = $1`, [channel.id]);

            for (const command of commands) {
                const role = roles.find((r: any) => r.type === command.permission);
                if (role) {
                    await queryRunner.query(`update "channels_commands" set "rolesIds" = array_append("rolesIds", $1) where "id" = $2`, [role.id, command.id]);
                }
            }
        }
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_roles_users" DROP CONSTRAINT "FK_c3b3d3917f24ec3fe3f5049e667"`);
        await queryRunner.query(`ALTER TABLE "channels_roles_users" DROP CONSTRAINT "FK_ae2e54203d7fed7ed38243e8b13"`);
        await queryRunner.query(`ALTER TABLE "channels_roles" DROP CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103"`);
        await queryRunner.query(`DROP TABLE "channels_roles_users"`);
        await queryRunner.query(`DROP TABLE "channels_roles"`);
        await queryRunner.query(`DROP TYPE "public"."channels_roles_permissions_enum"`);
        await queryRunner.query(`DROP TYPE "public"."channels_roles_type_enum"`);
    }

}
