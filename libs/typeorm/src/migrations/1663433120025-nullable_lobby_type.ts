import { MigrationInterface, QueryRunner } from "typeorm";

export class nullableLobbyType1663433120025 implements MigrationInterface {
    name = 'nullableLobbyType1663433120025'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "dota_matches" ALTER COLUMN "lobby_type" DROP NOT NULL`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "dota_matches" ALTER COLUMN "lobby_type" SET NOT NULL`);
    }

}
