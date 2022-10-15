import { MigrationInterface, QueryRunner } from "typeorm";

export class addProcessedToGreetings1665855111524 implements MigrationInterface {
    name = 'addProcessedToGreetings1665855111524'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_greetings" ADD "processed" boolean NOT NULL DEFAULT false`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_greetings" DROP COLUMN "processed"`);
    }

}
