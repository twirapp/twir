import { ArrayNotEmpty, IsArray, IsBoolean, IsNotEmpty, IsNumber, IsOptional, IsString, MaxLength, Min } from 'class-validator';

export class CreateTimerDto {
  @IsString()
  @IsNotEmpty()
  name: string;

  @IsBoolean()
  @IsOptional()
  enabled: boolean;

  @IsNumber()
  @Min(5)
  @IsNotEmpty()
  timeInterval: number;

  @IsNumber()
  @IsNotEmpty()
  messageInterval: number;

  @IsArray({})
  @ArrayNotEmpty()
  @IsString({ each: true })
  @IsNotEmpty({ each: true })
  @MaxLength(400, { each: true })
  responses: string[];
}
