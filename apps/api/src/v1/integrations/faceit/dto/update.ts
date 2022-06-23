
import { Type } from 'class-transformer';
import { IsBoolean, IsOptional, IsString, ValidateNested } from 'class-validator';

class Data {
  @IsString()
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