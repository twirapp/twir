import { MigrationInterface, QueryRunner } from 'typeorm';

export class commandsAddOnlineOnly1678880601877 implements MigrationInterface {
    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD "online_only" boolean NOT NULL DEFAULT false`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "online_only"`);
    }
}
