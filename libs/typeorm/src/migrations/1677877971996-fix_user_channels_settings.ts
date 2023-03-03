import { MigrationInterface, QueryRunner } from "typeorm";

export class fixUserChannelsSettings1677877971996 implements MigrationInterface {
    name = 'fixUserChannelsSettings1677877971996'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "FK_bc603823f3f741359c2339389f9"`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72"`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" DROP CONSTRAINT "FK_c17a48a983ac20ff800f553630a"`);
        await queryRunner.query(`ALTER TABLE "channels_permits" DROP CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06"`);
        await queryRunner.query(`ALTER TABLE "channel_events_list" DROP CONSTRAINT "FK_62386183f7a7575141594880b03"`);
        await queryRunner.query(`ALTER TABLE "channels_timers" DROP CONSTRAINT "FK_50978b1848b99458d6ceb6b1989"`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" DROP CONSTRAINT "FK_15402461d2eb224e2e088b35606"`);
        await queryRunner.query(`ALTER TABLE "users_online" DROP CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4"`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" DROP CONSTRAINT "FK_b34073350d4aa307b6104380e9b"`);
        await queryRunner.query(`ALTER TABLE "channels_messages" DROP CONSTRAINT "FK_05a589d7e48f170714dc73243bf"`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" DROP CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c"`);
        await queryRunner.query(`ALTER TABLE "channels_streams" DROP CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa"`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" DROP CONSTRAINT "FK_a757e3014566676c024e4ce16d1"`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" DROP CONSTRAINT "FK_309ade49a31238d00065fc7c32e"`);
        await queryRunner.query(`ALTER TABLE "users_stats" DROP CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258"`);
        await queryRunner.query(`ALTER TABLE "channels_info_history" DROP CONSTRAINT "FK_d326d1f33afefc8f45d6d546917"`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" DROP CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4"`);
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" DROP CONSTRAINT "FK_c202e2ed66394a1bd6651734078"`);
        await queryRunner.query(`ALTER TABLE "channels_roles" DROP CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "FK_c145b2745bd936041f37b5d5d49"`);
        await queryRunner.query(`ALTER TABLE "channels_events" DROP CONSTRAINT "FK_763ec88e86ecbf8ca6ec3a9ec7b"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_2ba87fd85e8e748470257452227"`);
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "UQ_bc603823f3f741359c2339389f9"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "UQ_e74a3ef66bba62b18e3448211f7" UNIQUE ("channelId", "userId")`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" ADD CONSTRAINT "FK_309ade49a31238d00065fc7c32e" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "FK_c145b2745bd936041f37b5d5d49" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_permits" ADD CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_roles" ADD CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" ADD CONSTRAINT "FK_c202e2ed66394a1bd6651734078" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_2ba87fd85e8e748470257452227" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "users_online" ADD CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "users_stats" ADD CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_messages" ADD CONSTRAINT "FK_05a589d7e48f170714dc73243bf" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ADD CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channel_events_list" ADD CONSTRAINT "FK_62386183f7a7575141594880b03" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" ADD CONSTRAINT "FK_15402461d2eb224e2e088b35606" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_info_history" ADD CONSTRAINT "FK_d326d1f33afefc8f45d6d546917" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" ADD CONSTRAINT "FK_c17a48a983ac20ff800f553630a" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" ADD CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "FK_b34073350d4aa307b6104380e9b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_streams" ADD CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa" FOREIGN KEY ("userId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_timers" ADD CONSTRAINT "FK_50978b1848b99458d6ceb6b1989" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_events" ADD CONSTRAINT "FK_763ec88e86ecbf8ca6ec3a9ec7b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "FK_bc603823f3f741359c2339389f9" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ADD CONSTRAINT "FK_a757e3014566676c024e4ce16d1" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" DROP CONSTRAINT "FK_a757e3014566676c024e4ce16d1"`);
        await queryRunner.query(`ALTER TABLE "channels" DROP CONSTRAINT "FK_bc603823f3f741359c2339389f9"`);
        await queryRunner.query(`ALTER TABLE "channels_events" DROP CONSTRAINT "FK_763ec88e86ecbf8ca6ec3a9ec7b"`);
        await queryRunner.query(`ALTER TABLE "channels_timers" DROP CONSTRAINT "FK_50978b1848b99458d6ceb6b1989"`);
        await queryRunner.query(`ALTER TABLE "channels_streams" DROP CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa"`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" DROP CONSTRAINT "FK_b34073350d4aa307b6104380e9b"`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" DROP CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c"`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" DROP CONSTRAINT "FK_c17a48a983ac20ff800f553630a"`);
        await queryRunner.query(`ALTER TABLE "channels_info_history" DROP CONSTRAINT "FK_d326d1f33afefc8f45d6d546917"`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" DROP CONSTRAINT "FK_15402461d2eb224e2e088b35606"`);
        await queryRunner.query(`ALTER TABLE "channel_events_list" DROP CONSTRAINT "FK_62386183f7a7575141594880b03"`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" DROP CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72"`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" DROP CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4"`);
        await queryRunner.query(`ALTER TABLE "channels_messages" DROP CONSTRAINT "FK_05a589d7e48f170714dc73243bf"`);
        await queryRunner.query(`ALTER TABLE "users_stats" DROP CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258"`);
        await queryRunner.query(`ALTER TABLE "users_online" DROP CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4"`);
        await queryRunner.query(`ALTER TABLE "channels_commands" DROP CONSTRAINT "FK_2ba87fd85e8e748470257452227"`);
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" DROP CONSTRAINT "FK_c202e2ed66394a1bd6651734078"`);
        await queryRunner.query(`ALTER TABLE "channels_roles" DROP CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103"`);
        await queryRunner.query(`ALTER TABLE "channels_permits" DROP CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "FK_c145b2745bd936041f37b5d5d49"`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" DROP CONSTRAINT "FK_309ade49a31238d00065fc7c32e"`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" DROP CONSTRAINT "UQ_e74a3ef66bba62b18e3448211f7"`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "UQ_bc603823f3f741359c2339389f9" UNIQUE ("id")`);
        await queryRunner.query(`ALTER TABLE "channels_commands" ADD CONSTRAINT "FK_2ba87fd85e8e748470257452227" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_events" ADD CONSTRAINT "FK_763ec88e86ecbf8ca6ec3a9ec7b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "FK_c145b2745bd936041f37b5d5d49" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_roles" ADD CONSTRAINT "FK_1c6f5f58e54b77d7480a4895103" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_commands_groups" ADD CONSTRAINT "FK_c202e2ed66394a1bd6651734078" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_customvars" ADD CONSTRAINT "FK_2bf1744a6cd76e4457eddfd3bc4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_info_history" ADD CONSTRAINT "FK_d326d1f33afefc8f45d6d546917" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "users_stats" ADD CONSTRAINT "FK_d55aab4a64c0c6b4b374b1da258" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_emotes_usages" ADD CONSTRAINT "FK_309ade49a31238d00065fc7c32e" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ADD CONSTRAINT "FK_a757e3014566676c024e4ce16d1" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_streams" ADD CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa" FOREIGN KEY ("userId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_keywords" ADD CONSTRAINT "FK_d17ea8fc949494c6b9681d5350c" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_messages" ADD CONSTRAINT "FK_05a589d7e48f170714dc73243bf" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_moderation_settings" ADD CONSTRAINT "FK_b34073350d4aa307b6104380e9b" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "users_online" ADD CONSTRAINT "FK_e6ae29713ab794b6ad8ef4fe5b4" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_greetings" ADD CONSTRAINT "FK_15402461d2eb224e2e088b35606" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_timers" ADD CONSTRAINT "FK_50978b1848b99458d6ceb6b1989" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channel_events_list" ADD CONSTRAINT "FK_62386183f7a7575141594880b03" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
        await queryRunner.query(`ALTER TABLE "channels_permits" ADD CONSTRAINT "FK_b7e83e12e0482fa968de97b6b06" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_integrations" ADD CONSTRAINT "FK_c17a48a983ac20ff800f553630a" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels_dota_accounts" ADD CONSTRAINT "FK_a4a2cd666dac0cae74549a5de72" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
        await queryRunner.query(`ALTER TABLE "channels" ADD CONSTRAINT "FK_bc603823f3f741359c2339389f9" FOREIGN KEY ("id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE`);
    }

}
