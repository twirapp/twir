import { MigrationInterface, QueryRunner } from 'typeorm';

export class ytsrAddDisplayName1672927752886 implements MigrationInterface {
    name = 'ytsrAddDisplayName1672927752886';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ADD "orderedByDisplayName" character varying`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ALTER COLUMN "id" DROP DEFAULT`);
    }

}
