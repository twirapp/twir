import { MigrationInterface, QueryRunner } from "typeorm";

export class cascadeTimersResponses1666371747377 implements MigrationInterface {
    name = 'cascadeTimersResponses1666371747377'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_timers_responses" DROP CONSTRAINT "FK_464cac4218d2ba9ac84325d06c8"`);
        await queryRunner.query(`ALTER TABLE "channels_timers_responses" ADD CONSTRAINT "FK_464cac4218d2ba9ac84325d06c8" FOREIGN KEY ("timerId") REFERENCES "channels_timers"("id") ON DELETE CASCADE ON UPDATE NO ACTION`);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`ALTER TABLE "channels_timers_responses" DROP CONSTRAINT "FK_464cac4218d2ba9ac84325d06c8"`);
        await queryRunner.query(`ALTER TABLE "channels_timers_responses" ADD CONSTRAINT "FK_464cac4218d2ba9ac84325d06c8" FOREIGN KEY ("timerId") REFERENCES "channels_timers"("id") ON DELETE NO ACTION ON UPDATE NO ACTION`);
    }

}
