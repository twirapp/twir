import { MigrationInterface, QueryRunner } from 'typeorm';

export class addChannelInfoHistory1677006049732 implements MigrationInterface {
    name = 'addChannelInfoHistory1677006049732';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "channels_info_history" ("id" uuid NOT NULL DEFAULT uuid_generate_v4(), "channelId" text NOT NULL, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), "title" text NOT NULL, "category" character varying NOT NULL, CONSTRAINT "PK_deac5bf159894268f422f66c168" PRIMARY KEY ("id"))`);
        await queryRunner.query(`ALTER TABLE "channels_info_history" ADD CONSTRAINT "FK_d326d1f33afefc8f45d6d546917" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_info_history" DROP CONSTRAINT "FK_d326d1f33afefc8f45d6d546917"`);
        await queryRunner.query(`DROP TABLE "channels_info_history"`);
    }

}
