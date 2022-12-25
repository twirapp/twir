import { MigrationInterface, QueryRunner } from 'typeorm';

export class nonNullableKeywordsCooldown1671194430883 implements MigrationInterface {
    name = 'nonNullableKeywordsCooldown1671194430883';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_keywords" ALTER COLUMN "cooldown" SET NOT NULL`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_keywords" ALTER COLUMN "cooldown" DROP NOT NULL`);
    }

}
