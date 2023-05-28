import { MigrationInterface, QueryRunner } from 'typeorm';

export class addCommandsRestrictions1685295875918 implements MigrationInterface {
	name = 'addCommandsRestrictions1685295875918';

	public async up(queryRunner: QueryRunner): Promise<void> {
		await queryRunner.query(
			`ALTER TABLE "channels_commands" ADD "requiredWatchTime" integer NOT NULL DEFAULT '0'`,
		);
		await queryRunner.query(
			`ALTER TABLE "channels_commands" ADD "requiredMessages" integer NOT NULL DEFAULT '0'`,
		);
		await queryRunner.query(
			`ALTER TABLE "channels_commands" ADD "requiredUsedChannelPoints" integer NOT NULL DEFAULT '0'`,
		);
	}

	public async down(queryRunner: QueryRunner): Promise<void> {
		await queryRunner.query(
			`ALTER TABLE "channels_commands" DROP COLUMN "requiredUsedChannelPoints"`,
		);
		await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "requiredMessages"`);
		await queryRunner.query(`ALTER TABLE "channels_commands" DROP COLUMN "requiredWatchTime"`);
	}
}
