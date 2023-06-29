import {
  ChannelModerationSetting,
  SettingsType,
} from '@twir/typeorm/entities/ChannelModerationSetting';
import {
  IsArray,
  IsBoolean,
  IsIn,
  IsNumber,
  IsOptional,
  IsString,
  ValidateNested,
} from 'class-validator';

export class ModerationSettingsDto implements Omit<ChannelModerationSetting, 'id' | 'channelId'> {
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
  blackListSentences: any[];
}

export class ModerationUpdateDto {
  @ValidateNested({ each: true })
  items: ModerationSettingsDto[];
}
