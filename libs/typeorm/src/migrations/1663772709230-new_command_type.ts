import { MigrationInterface, QueryRunner } from 'typeorm';

export class newCommandType1663772709230 implements MigrationInterface {
  name = 'newCommandType1663772709230';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TYPE "public"."channels_commands_module_enum" RENAME TO "channels_commands_module_enum_old"`,
    );
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_module_enum" AS ENUM('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION', 'MANAGE')`,
    );
    await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum" USING "module"::"text"::"public"."channels_commands_module_enum"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum_old"`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `CREATE TYPE "public"."channels_commands_module_enum_old" AS ENUM('CUSTOM', 'DOTA', 'CHANNEL', 'MODERATION')`,
    );
    await queryRunner.query(`ALTER TABLE "channels_commands" ALTER COLUMN "module" DROP DEFAULT`);
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" TYPE "public"."channels_commands_module_enum_old" USING "module"::"text"::"public"."channels_commands_module_enum_old"`,
    );
    await queryRunner.query(
      `ALTER TABLE "channels_commands" ALTER COLUMN "module" SET DEFAULT 'CUSTOM'`,
    );
    await queryRunner.query(`DROP TYPE "public"."channels_commands_module_enum"`);
    await queryRunner.query(
      `ALTER TYPE "public"."channels_commands_module_enum_old" RENAME TO "channels_commands_module_enum"`,
    );
  }
}
