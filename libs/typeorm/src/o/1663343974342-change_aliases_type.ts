import { MigrationInterface, QueryRunner, TableColumn } from 'typeorm';

export class changeAliasesType1663343974342 implements MigrationInterface {
  name = 'changeAliasesType1663343974342';

  public async up(queryRunner: QueryRunner): Promise<void> {
    const column = new TableColumn({
      name: 'newAliases',
      type: 'text',
      isArray: true,
      default: `'{}'`,
    });
    await queryRunner.addColumn('channels_commands', column);

    const commands = await queryRunner.query(`SELECT * from "channels_commands"`);
    for (const command of commands) {
      await queryRunner.query(`UPDATE "channels_commands" SET "newAliases"=$1 WHERE "id"=$2`, [
        command.aliases,
        command.id,
      ]);
    }

    await queryRunner.dropColumn('channels_commands', 'aliases');
    await queryRunner.renameColumn('channels_commands', 'newAliases', 'aliases');
  }

  public async down(_queryRunner: QueryRunner): Promise<void> {
    return;
  }
}
