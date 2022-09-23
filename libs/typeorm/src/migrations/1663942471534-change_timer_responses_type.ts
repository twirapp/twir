import { MigrationInterface, QueryRunner, TableColumn } from 'typeorm';

export class changeTimerResponsesType1663942471534 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    const column = new TableColumn({
      name: 'newResponses',
      type: 'text',
      isArray: true,
      default: `'{}'`,
    });
    await queryRunner.addColumn('channels_timers', column);

    const timers = await queryRunner.query(`SELECT * from "channels_timers"`);
    for (const timer of timers) {
      await queryRunner.query(`UPDATE "channels_timers" SET "newResponses"=$1 WHERE "id"=$2`, [
        timer.responses,
        timer.id,
      ]);
    }

    await queryRunner.dropColumn('channels_timers', 'responses');
    await queryRunner.renameColumn('channels_timers', 'newResponses', 'responses');
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    return;
  }
}
