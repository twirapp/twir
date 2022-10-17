import { MigrationInterface, QueryRunner } from "typeorm";

export class changeStreamStartedAtType1666008095838 implements MigrationInterface {
    name = 'changeStreamStartedAtType1666008095838'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_streams" DROP COLUMN "startedAt"`);
        await queryRunner.query(`ALTER TABLE "channels_streams" ADD "startedAt" TIMESTAMP NOT NULL`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_streams" DROP COLUMN "startedAt"`);
        await queryRunner.query(`ALTER TABLE "channels_streams" ADD "startedAt" text NOT NULL`);
    }

}
