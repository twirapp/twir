import { MigrationInterface, QueryRunner } from 'typeorm';

export class addSongLinkToRequests1686693093611 implements MigrationInterface {
    name = 'addSongLinkToRequests1686693093611';

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" ADD "songLink" character varying`);

				const settingsRows = await queryRunner.query(
					`SELECT * from "channels_modules_settings" WHERE "type" = $1`,
					['youtube_song_requests'],
				);

				for (const row of settingsRows) {
					await queryRunner.query(
						`UPDATE "channels_modules_settings" SET "settings" = $1 WHERE "id" = $2`,
						[{
							...row.settings,
							translations: {
								...row.settings.translations,
								nowPlaying: 'Now playing "{{songTitle}}" {{songLink}} requested by @{{orderedByDisplayName}}',
							},
						}, row.id],
					);
				}
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_requested_songs" DROP COLUMN "songLink"`);
    }

}
