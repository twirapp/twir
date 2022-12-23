import { MigrationInterface, QueryRunner } from 'typeorm';

export class nonNullableVarResponseAndEval1671194276298 implements MigrationInterface {
    name = 'nonNullableVarResponseAndEval1671194276298';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "evalValue" SET DEFAULT ''`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "response" SET DEFAULT ''`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "evalValue" SET NOT NULL`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "response" SET NOT NULL`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "response" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "evalValue" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "response" DROP NOT NULL`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ALTER COLUMN "evalValue" DROP NOT NULL`);
    }

}
