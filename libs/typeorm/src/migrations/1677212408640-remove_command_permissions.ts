import { MigrationInterface, QueryRunner } from "typeorm";

export class removeCommandPermissions1677212408640 implements MigrationInterface {
    name = 'removeCommandPermissions1677212408640'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "permission"`);
        await queryRunner.query(`DROP TYPE "public"."channels_commands_permission_enum"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."channels_commands_permission_enum" AS ENUM('BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER')`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "permission" "public"."channels_commands_permission_enum" NOT NULL`);
    }

}
