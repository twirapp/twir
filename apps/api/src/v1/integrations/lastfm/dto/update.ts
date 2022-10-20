import { Type } from 'class-transformer';
import {
  IsBoolean,
  IsNotEmpty,
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
}

export class LastfmUpdateDto {
  @ValidateNested()
  @ValidateIf((o: LastfmUpdateDto) => o.enabled)
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}
