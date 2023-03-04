import { MigrationInterface, QueryRunner } from 'typeorm';

export class fixUserChannelsSettings21677923559083 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "UQ_b5f1c883e497ba7a0eeae08e8b8"`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
