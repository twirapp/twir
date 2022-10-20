import { Type } from 'class-transformer';
import {
  IsBoolean,
  IsNotEmpty,
  IsOptional,
  IsString,
  MaxLength,
  ValidateIf,
  ValidateNested,
} from 'class-validator';

class Data {
  @IsString()
  @IsNotEmpty()
  @MaxLength(100)
  username: string;

  @IsString()
  @IsOptional()
  game?: string;
}

export class FaceitUpdateDto {
  @ValidateNested()
  @ValidateIf((o: FaceitUpdateDto) => o.enabled)
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}
