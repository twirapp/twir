import { MigrationInterface, QueryRunner } from 'typeorm';

export class dropIntegrations1671823591036 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`DELETE FROM channels_integrations`);
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-function
    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
