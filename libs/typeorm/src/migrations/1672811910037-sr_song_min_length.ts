import { MigrationInterface, QueryRunner } from 'typeorm';

export class srSongMinLength1672811910037 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        const currentSettings = await queryRunner.query('SELECT * from "channels_modules_settings"');

        for (const row of currentSettings) {
            const settings = row.settings;

            await queryRunner.query(
              'UPDATE "channels_modules_settings" SET "settings"=$1 WHERE "channelId"=$2',
              [
                  {
                      ...settings,
                      song: {
                        ...settings.song,
                        minLength: 0,
                      },
                  },
                  row.channelId,
              ],
            );
        }
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
