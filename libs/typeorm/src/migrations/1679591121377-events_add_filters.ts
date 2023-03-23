import { MigrationInterface, QueryRunner } from 'typeorm';

export class eventsAddFilters1679591121377 implements MigrationInterface {
    name = 'eventsAddFilters1679591121377';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TYPE "public"."channels_events_operations_filters_type_enum" AS ENUM('EQUALS', 'NOT_EQUALS', 'CONTAINS', 'NOT_CONTAINS', 'STARTS_WITH', 'ENDS_WITH', 'GREATER_THAN', 'LESS_THAN', 'GREATER_THAN_OR_EQUALS', 'LESS_THAN_OR_EQUALS', 'IS_EMPTY', 'IS_NOT_EMPTY')`);
        await queryRunner.query(`CREATE TABLE "channels_events_operations_filters" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "operationId" uuid NOT NULL, "type" "public"."channels_events_operations_filters_type_enum" NOT NULL, "left" text NOT NULL, "right" text NOT NULL, CONSTRAINT "PK_184c3468090c81affe77014f162" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_events_operations_filters" ADD CONSTRAINT "FK_d56a4ca65d44fc4adbaf66a2e80" FOREIGN KEY ("operationId") REFERENCES "channels_events_operations"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_events_operations_filters" DROP CONSTRAINT "FK_d56a4ca65d44fc4adbaf66a2e80"`);
        await queryRunner.query(`DROP TABLE "channels_events_operations_filters"`);
        await queryRunner.query(`DROP TYPE "public"."channels_events_operations_filters_type_enum"`);
    }

}
