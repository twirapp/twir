import { CommandPermission, Response, CooldownType } from '@tsuwari/prisma';
import { ArrayNotEmpty, IsArray, IsBoolean, IsIn, IsNotEmpty, IsNumber, IsOptional, IsString, MaxLength } from 'class-validator';
import type { SetOptional } from 'type-fest';

export class UpdateOrCreateCommandDto {
  id?: string;

  @IsString()
  @IsNotEmpty()
  name: string;

  @IsNumber()
  @IsOptional()
  cooldown?: number;

  @IsIn(Object.keys(CooldownType))
  @IsOptional()
  cooldownType?: CooldownType;

  @IsString()
  @IsOptional()
  description?: string | null;

  @IsIn(Object.keys(CommandPermission))
  permission: CommandPermission;

  @IsArray()
  @IsOptional()
  aliases?: string[];

  @IsBoolean()
  @IsOptional()
  visible?: boolean;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;

  @IsArray()
  @ArrayNotEmpty()
  @MaxLength(400, { each: true })
  responses: Array<SetOptional<Omit<Response, 'commandId'>, 'id'>>;
}
