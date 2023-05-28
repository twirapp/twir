import { MigrationInterface, QueryRunner } from "typeorm";

export class AddCategoryIdForAlias1685245990972 implements MigrationInterface {
    name = 'AddCategoryIdForAlias1685245990972'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_categories_aliases" ADD "categoryId" text NOT NULL`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_categories_aliases" DROP COLUMN "categoryId"`);
    }

}
