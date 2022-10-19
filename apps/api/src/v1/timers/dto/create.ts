import { Type } from 'class-transformer';
import {
  ArrayNotEmpty,
  IsArray,
  IsBoolean,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsString,
  Max,
  MaxLength,
  Min,
  ValidateNested,
} from 'class-validator';

export class CreateTimerDto {
  @IsString()
  @IsNotEmpty()
  @MaxLength(510)
  name: string;

  @IsBoolean()
  @IsOptional()
  enabled: boolean;

  @IsNumber()
  @Min(1)
  @Max(120)
  @IsNotEmpty()
  timeInterval: number;

  @IsNumber()
  @IsNotEmpty()
  @Max(10000)
  messageInterval: number;

  @IsArray()
  @ValidateNested()
  @Type(() => TimerResponse)
  @ArrayNotEmpty()
  responses: Array<TimerResponse>;
}

class TimerResponse {
  @IsString()
  @IsNotEmpty()
  @MaxLength(400)
  text: string;

  @IsBoolean()
  @IsOptional()
  isAnnounce?: boolean;
}
