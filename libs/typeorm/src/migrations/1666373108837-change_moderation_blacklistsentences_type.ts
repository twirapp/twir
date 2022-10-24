import { MigrationInterface, QueryRunner, TableColumn } from 'typeorm';

export class changeModerationBlacklistsentencesType1666373108837 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    const column = new TableColumn({
      name: 'newBlackListSentences',
      type: 'text',
      isArray: true,
      default: `'{}'`,
    });
    await queryRunner.addColumn('channels_moderation_settings', column);

    const settings = await queryRunner.query(`SELECT * from "channels_moderation_settings"`);
    for (const setting of settings) {
      await queryRunner.query(
        `UPDATE "channels_moderation_settings" SET "newBlackListSentences"=$1 WHERE "id"=$2`,
        [setting.blackListSentences, setting.id],
      );
    }

    await queryRunner.dropColumn('channels_moderation_settings', 'blackListSentences');
    await queryRunner.renameColumn(
      'channels_moderation_settings',
      'newBlackListSentences',
      'blackListSentences',
    );
  }

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  public async down(queryRunner: QueryRunner): Promise<void> {}
}
