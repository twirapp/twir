import { MigrationInterface, QueryRunner } from "typeorm"

export class renameSongrequestsToSongs1680944512410 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TYPE channels_commands_module_enum RENAME VALUE 'SONGREQUEST' TO 'SONGS'`)
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
