import { MigrationInterface, QueryRunner } from 'typeorm';

export class dropDonatelloChannelsIntegrations1685541244421 implements MigrationInterface {
	public async up(queryRunner: QueryRunner): Promise<void> {
		const integrations = await queryRunner.query(
			`SELECT * FROM "integrations" WHERE "service" = $1`,
			['DONATELLO'],
		);

		const integration = integrations.at(0);

		if (!integration) {
			throw new Error('Integration not found');
		}

		await queryRunner.query(`DELETE FROM "channels_integrations" WHERE "integrationId" = $1`, [
			integration.id,
		]);
	}

	public async down(queryRunner: QueryRunner): Promise<void> {}
}
