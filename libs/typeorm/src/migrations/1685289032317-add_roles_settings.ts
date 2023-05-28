import { MigrationInterface, QueryRunner } from 'typeorm';

export class addRolesSettings1685289032317 implements MigrationInterface {
	name = 'addRolesSettings1685289032317';

	public async up(queryRunner: QueryRunner): Promise<void> {
		await queryRunner.query(
			`ALTER TABLE "channels_roles" ADD "settings" jsonb NOT NULL DEFAULT '{}'`,
		);
	}

	public async down(queryRunner: QueryRunner): Promise<void> {
		await queryRunner.query(`ALTER TABLE "channels_roles" DROP COLUMN "settings"`);
	}
}
