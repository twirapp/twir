
import { Type } from 'class-transformer';
import { IsBoolean, IsNotEmpty, IsOptional, IsString, ValidateNested } from 'class-validator';

class Data {
  @IsString()
  @IsNotEmpty()
  username: string;

  @IsString()
  @IsOptional()
  game?: string;
}

export class FaceitUpdateDto {
  @ValidateNested()
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}