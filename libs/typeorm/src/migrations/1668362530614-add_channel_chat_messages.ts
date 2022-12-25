import { MigrationInterface, QueryRunner } from "typeorm";

export class addChannelChatMessages1668362530614 implements MigrationInterface {
    name = 'addChannelChatMessages1668362530614'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`CREATE TABLE "channels_messages" ("messageId" text NOT NULL, "channelId" text NOT NULL, "userId" text NOT NULL, "userName" text NOT NULL, "text" text NOT NULL, "canBeDeleted" boolean NOT NULL DEFAULT true, "createdAt" TIMESTAMP NOT NULL DEFAULT now(), CONSTRAINT "PK_685840c6efdbd345cf265976753" PRIMARY KEY ("messageId"))`);
        await queryRunner.query(`CREATE INDEX "IDX_05a589d7e48f170714dc73243b" ON "channels_messages" ("channelId") `);
        await queryRunner.query(`CREATE INDEX "IDX_dd5560724c1166b9b70954f44b" ON "channels_messages" ("userId") `);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" ALTER COLUMN "blackListSentences" DROP NOT NULL`);
        await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" SET DEFAULT gen_random_uuid()`);
        await queryRunner.query(`ALTER TABLE "channels_messages" ADD CONSTRAINT "FK_05a589d7e48f170714dc73243bf" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_messages" ADD CONSTRAINT "FK_dd5560724c1166b9b70954f44be" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_messages" DROP CONSTRAINT "FK_dd5560724c1166b9b70954f44be"`);
        await queryRunner.query(`ALTER TABLE "channels_messages" DROP CONSTRAINT "FK_05a589d7e48f170714dc73243bf"`);
        await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" DROP DEFAULT`);
        await queryRunner.query(`ALTER TABLE "users" ALTER COLUMN "apiKey" SET DEFAULT uuid_generate_v4()`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" ALTER COLUMN "blackListSentences" SET NOT NULL`);
        await queryRunner.query(`DROP INDEX "public"."IDX_dd5560724c1166b9b70954f44b"`);
        await queryRunner.query(`DROP INDEX "public"."IDX_05a589d7e48f170714dc73243b"`);
        await queryRunner.query(`DROP TABLE "channels_messages"`);
    }

}
