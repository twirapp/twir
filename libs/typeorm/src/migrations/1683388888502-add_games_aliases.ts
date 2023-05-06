import { MigrationInterface, QueryRunner } from "typeorm";

export class AddGamesAliases1683388888502 implements MigrationInterface {
    name = 'AddGamesAliases1683388888502'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "channels_categories_aliases" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "category" text NOT NULL, "alias" text NOT NULL, "channelId" text NOT NULL, CONSTRAINT "UQ_c5f8b666ae88b282bb3e543375f" UNIQUE ("alias", "channelId"), CONSTRAINT "PK_feee17b6a79b9238e6898142c88" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE INDEX "IDX_48bc826999f991293f9029af26" ON "channels_categories_aliases" ("channelId") `);
        await queryRunner.query(`ALTER TABLE "channels_categories_aliases" ADD CONSTRAINT "FK_48bc826999f991293f9029af26a" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_categories_aliases" DROP CONSTRAINT "FK_48bc826999f991293f9029af26a"`);
        await queryRunner.query(`DROP INDEX "public"."IDX_48bc826999f991293f9029af26"`);
        await queryRunner.query(`DROP TABLE "channels_categories_aliases"`);
    }

}
