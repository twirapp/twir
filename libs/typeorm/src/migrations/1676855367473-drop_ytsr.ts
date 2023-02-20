import { MigrationInterface, QueryRunner } from 'typeorm';

export class dropYtsr1676855367473 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`DELETE FROM "channels_commands" WHERE "defaultName" = $1`, ['ytsr']);
        await queryRunner.query(`DELETE FROM "channels_commands" WHERE "defaultName" = $1`, ['ytsr wrong']);
        await queryRunner.query(`DELETE FROM "channels_commands" WHERE "defaultName" = $1`, ['ytsr list']);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
