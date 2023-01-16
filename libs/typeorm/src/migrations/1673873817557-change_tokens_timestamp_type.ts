import { MigrationInterface, QueryRunner } from 'typeorm';

export class changeTokensTimestampType1673873817557 implements MigrationInterface {
    name = 'changeTokensTimestampType1673873817557';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "tokens" ALTER COLUMN "obtainmentTimestamp" TYPE TIMESTAMP WITH TIME ZONE`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "tokens" ALTER COLUMN "obtainmentTimestamp" TYPE TIMESTAMP`);
    }

}
