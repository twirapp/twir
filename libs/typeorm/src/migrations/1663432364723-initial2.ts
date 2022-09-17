import { MigrationInterface, QueryRunner } from "typeorm";

export class initial21663432364723 implements MigrationInterface {
    name = 'initial21663432364723'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "FK_bc603823f3f741359c2339389f9"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_2ba87fd85e8e748470257452227"`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" DROP CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4"`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72"`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" DROP CONSTRAINT "FK_15402461d2eb224e2e088b35606"`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" DROP CONSTRAINT "FK_c17a48a983ac20ff800f553630a"`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" DROP CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c"`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" DROP CONSTRAINT "FK_b34073350d4aa307b6104380e9b"`);
        await queryRunner.query(`ALTER TABLE "channels_permits" DROP CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06"`);
        await queryRunner.query(`ALTER TABLE "channels_timers" DROP CONSTRAINT "FK_50978b1848b99458d6ceb6b1989"`);
        await queryRunner.query(`ALTER TABLE "channels_dashboard_access" DROP CONSTRAINT "FK_82931b7f57c891fd7da375f3892"`);
        await queryRunner.query(`ALTER TABLE "users_stats" DROP CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258"`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "UQ_bc603823f3f741359c2339389f9" UNIQUE ("id")`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "FK_bc603823f3f741359c2339389f9" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_2ba87fd85e8e748470257452227" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ADD CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" ADD CONSTRAINT "FK_15402461d2eb224e2e088b35606" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" ADD CONSTRAINT "FK_c17a48a983ac20ff800f553630a" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" ADD CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "FK_b34073350d4aa307b6104380e9b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_permits" ADD CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_timers" ADD CONSTRAINT "FK_50978b1848b99458d6ceb6b1989" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "FK_82931b7f57c891fd7da375f3892" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "users_stats" ADD CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "users_stats" DROP CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258"`);
        await queryRunner.query(`ALTER TABLE "channels_dashboard_access" DROP CONSTRAINT "FK_82931b7f57c891fd7da375f3892"`);
        await queryRunner.query(`ALTER TABLE "channels_timers" DROP CONSTRAINT "FK_50978b1848b99458d6ceb6b1989"`);
        await queryRunner.query(`ALTER TABLE "channels_permits" DROP CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06"`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" DROP CONSTRAINT "FK_b34073350d4aa307b6104380e9b"`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" DROP CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c"`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" DROP CONSTRAINT "FK_c17a48a983ac20ff800f553630a"`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" DROP CONSTRAINT "FK_15402461d2eb224e2e088b35606"`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72"`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" DROP CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_2ba87fd85e8e748470257452227"`);
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "FK_bc603823f3f741359c2339389f9"`);
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "UQ_bc603823f3f741359c2339389f9"`);
        await queryRunner.query(`ALTER TABLE "users_stats" ADD CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_dashboard_access" ADD CONSTRAINT "FK_82931b7f57c891fd7da375f3892" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_timers" ADD CONSTRAINT "FK_50978b1848b99458d6ceb6b1989" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_permits" ADD CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "FK_b34073350d4aa307b6104380e9b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" ADD CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" ADD CONSTRAINT "FK_c17a48a983ac20ff800f553630a" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" ADD CONSTRAINT "FK_15402461d2eb224e2e088b35606" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ADD CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_2ba87fd85e8e748470257452227" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "FK_bc603823f3f741359c2339389f9" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
    }

}
