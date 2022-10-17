import { MigrationInterface, QueryRunner } from "typeorm";

export class makeKeywordCooldownNullable1666008755719 implements MigrationInterface {
    name = 'makeKeywordCooldownNullable1666008755719'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_keywords" ALTER COLUMN "cooldownExpireAt" DROP NOT NULL`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_keywords" ALTER COLUMN "cooldownExpireAt" SET NOT NULL`);
    }

}
