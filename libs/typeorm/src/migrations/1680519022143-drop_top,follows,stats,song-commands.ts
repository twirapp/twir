import { MigrationInterface, QueryRunner } from 'typeorm';

export class DropTopFollowsStatsSongCommands1680519022143 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    const commands = [
      'followage',
      'followsince',
      'song',
      'currentsong',
      'me',
      'top time',
      'top messages',
      'top watched',
      'top emotes',
      'top emotes users',
      'top points',
      'watchtime',
    ];

    for (const command of commands) {
      await queryRunner.query(
        `DELETE from "channels_commands" WHERE "name" = $1 AND "module" = 'CUSTOM'`,
        [command],
      );
    }
  }

  public async down(queryRunner: QueryRunner): Promise<void> {}
}
