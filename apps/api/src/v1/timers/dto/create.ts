import { IsArray, IsBoolean, IsNumber, IsOptional, IsString, Min } from 'class-validator';

export class CreateTimerDto {
  @IsString()
  name: string;

  @IsBoolean()
  @IsOptional()
  enabled: boolean;

  @IsNumber()
  @Min(5)
  timeInterval: number;

  @IsNumber()
  messageInterval: number;

  @IsArray({})
  responses: string[];
}
