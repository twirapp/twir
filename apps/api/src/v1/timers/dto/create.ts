import { ArrayNotEmpty, IsArray, IsBoolean, IsNotEmpty, IsNumber, IsOptional, IsString, Min } from 'class-validator';

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
  responses: string[];
}
