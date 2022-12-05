import { MigrationInterface, QueryRunner } from 'typeorm';

export class addCommandRestrictions1670198747795 implements MigrationInterface {
    name = 'addCommandRestrictions1670198747795';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."commands_restrictions_type_enum" AS ENUM('WATCHED', 'MESSAGES')`);
        await queryRunner.query(`CREATE TABLE "commands_restrictions" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "commandId" uuid NOT NULL, "type" "public"."commands_restrictions_type_enum" NOT NULL, "value" text NOT NULL, CONSTRAINT "PK_b78c9b42c1c870504bcac6f92f8" PRIMARY KEY ("id"))`);
        await queryRunner.query(`CREATE UNIQUE INDEX "channels_restrictions_commandId_type" ON "commands_restrictions" ("commandId", "type") `);
        await queryRunner.query(`ALTER TABLE "commands_restrictions" ADD CONSTRAINT "FK_f014b2dd4be25f9485a070059ea" FOREIGN KEY ("commandId") REFERENCES "channels_commands"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "commands_restrictions" DROP CONSTRAINT "FK_f014b2dd4be25f9485a070059ea"`);
        await queryRunner.query(`DROP INDEX "public"."channels_restrictions_commandId_type"`);
        await queryRunner.query(`DROP TABLE "commands_restrictions"`);
        await queryRunner.query(`DROP TYPE "public"."commands_restrictions_type_enum"`);
    }

}
