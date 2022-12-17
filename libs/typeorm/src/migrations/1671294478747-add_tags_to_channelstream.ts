import { MigrationInterface, QueryRunner } from "typeorm";

export class addTagsToChannelstream1671294478747 implements MigrationInterface {
    name = 'addTagsToChannelstream1671294478747'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_streams" ADD "tags" text array DEFAULT '{}'`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_streams" DROP COLUMN "tags"`);
    }

}
