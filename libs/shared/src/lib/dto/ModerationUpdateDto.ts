import { ModerationSettings, Prisma, SettingsType } from '@tsuwari/prisma';
import { IsArray, IsBoolean, IsIn, IsNumber, IsOptional, IsString, ValidateNested } from 'class-validator';

export class ModerationSettingsDto implements Omit<ModerationSettings, 'id' | 'channelId'> {
  @IsIn(Object.values(SettingsType))
  @IsString()
  type: SettingsType;

  @IsBoolean()
  @IsOptional()
  enabled: boolean;

  @IsBoolean()
  @IsOptional()
  subscribers: boolean;

  @IsBoolean()
  @IsOptional()
  vips: boolean;

  @IsNumber()
  @IsOptional()
  banTime: number;

  @IsString()
  @IsOptional()
  banMessage: string;

  @IsString()
  @IsOptional()
  warningMessage: string;

  @IsBoolean()
  @IsOptional()
  checkClips: boolean | null;

  @IsNumber()
  @IsOptional()
  triggerLength: number | null;

  @IsNumber()
  @IsOptional()
  maxPercentage: number | null;

  @IsArray()
  @IsOptional()
  @IsString({ each: true })
  blackListSentences: Prisma.JsonValue;
}

export class ModerationUpdateDto {
  @ValidateNested({ each: true })
  items: ModerationSettingsDto[];
}