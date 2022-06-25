import { ArrayNotEmpty, IsArray, IsBoolean, IsNotEmpty, IsNumber, IsOptional, IsString, Max, MaxLength, Min } from 'class-validator';

export class CreateTimerDto {
  @IsString()
  @IsNotEmpty()
  @MaxLength(510)
  name: string;

  @IsBoolean()
  @IsOptional()
  enabled: boolean;

  @IsNumber()
  @Min(5)
  @Max(84000)
  @IsNotEmpty()
  timeInterval: number;

  @IsNumber()
  @IsNotEmpty()
  @Max(10000)
  messageInterval: number;

  @IsArray()
  @ArrayNotEmpty()
  @IsString({ each: true })
  @IsNotEmpty({ each: true })
  @MaxLength(400, { each: true })
  responses: string[];
}
