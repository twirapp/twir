import { IsBoolean, IsNumber, IsOptional, IsString } from 'class-validator';


export class CreateKeywordDto {
  @IsString()
  text: string;

  @IsString()
  response: string;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;

  @IsNumber()
  @IsOptional()
  cooldown?: number;
}