import { IsBoolean, IsNumber, IsOptional, IsString, Min } from 'class-validator';


export class CreateKeywordDto {
  @IsString()
  text: string;

  @IsString()
  response: string;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;

  @IsNumber()
  @Min(5)
  @IsOptional()
  cooldown?: number;
}