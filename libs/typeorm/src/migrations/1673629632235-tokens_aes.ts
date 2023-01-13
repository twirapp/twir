import { config } from '@tsuwari/config';
import { encrypt } from '@tsuwari/crypto';
import { MigrationInterface, QueryRunner } from 'typeorm';

export class tokensAes1673629632235 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        const key = config.TOKENS_CIPHER_KEY;

        const tokens = await queryRunner.query('SELECT * from "tokens"');

        for (const t of tokens) {
            await queryRunner.query(
              'UPDATE "tokens" SET "accessToken"=$1, "refreshToken"=$2 WHERE "id"=$3',
              [
                  encrypt(t.accessToken, key),
                  encrypt(t.refreshToken, key),
                  t.id,
              ],
            );
        }
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
