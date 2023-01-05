import { MigrationInterface, QueryRunner } from 'typeorm';

export class ytsrTranslations1672908669216 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        const currentSettings = await queryRunner.query('SELECT * from "channels_modules_settings"');

        for (const row of currentSettings) {
            const settings = row.settings;

            await queryRunner.query(
              'UPDATE "channels_modules_settings" SET "settings"=$1 WHERE "channelId"=$2',
              [
                  {
                      ...settings,
                      translations: {
                        notEnabled: 'Song requests not enabled.',
                        nowPlaying: 'Now playing "{{title}} youtu.be/{{videoId}}" requested from @{{orderedByName}}',
                        noText: 'You should provide text for song request.',
                        acceptOnlyWhenOnline: 'Requests accepted only on online streams.',
                        song: {
                          notFound: 'Song not found.',
                          alreadyInQueue: 'Song already in queue.',
                          ageRestrictions: 'Age restriction on that song.',
                          cannotGetInformation: 'Cannot get information about song.',
                          live: 'Seems like that song is live, which is disallowed.',
                          denied: 'That song is denied for requests.',
                          notEnoughViews: 'Song haven\'t {{views}} views for request',
                        },
                        user: {
                          denied: 'You are denied to request any song.',
                        },
                        channel: {
                          denied: 'That channel is denied for requests.',
                        },
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
