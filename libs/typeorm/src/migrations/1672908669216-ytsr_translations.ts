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
                        nowPlaying: 'Now playing "{{songTitle}} youtu.be/{{songId}}" requested from @{{orderedByDisplayName}}',
                        noText: 'You should provide text for song request.',
                        acceptOnlyWhenOnline: 'Requests accepted only on online streams.',
                        song: {
                          notFound: 'Song not found.',
                          alreadyInQueue: 'Song already in queue.',
                          ageRestrictions: 'Age restriction on that song.',
                          cannotGetInformation: 'Cannot get information about song.',
                          live: 'Seems like that song is live, which is disallowed.',
                          denied: 'That song is denied for requests.',
                          requestedMessage: 'Song "{{songTitle}}" requested, queue position {{position}}. Estimated wait time before your track will be played is {{waitTime}}.',
                          maximumOrdered: 'Maximum number of songs is queued ({{maximum}}).',
                          minViews: 'Song {{songTitle}} ({{songViews}} views) haven\'t {{neededViews}} views for being ordered',
                          maxLength: 'Maximum length of song is {{maxLength}}',
                          minLength: 'Minimum length of song is {{minLength}}',
                        },
                        user: {
                          denied: 'You are denied to request any song.',
                          maxRequests: 'Maximum number of songs ordered by you ({{count}})',
                          minMessages: 'You have only {{userMessages}} messages, but needed {{neededMessages}} for requesting song',
                          minWatched: 'You have only {{userWatched}} messages, but needed {{neededWatched}} for requesting song',
                          minFollow: 'You are followed for {{userFollow}} minutes, but needed {{neededFollow}} for requesting song',
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
