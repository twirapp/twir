
import { Type } from 'class-transformer';
import { IsBoolean, IsString, ValidateNested } from 'class-validator';

class Data {
  @IsString()
  username: string;
}

export class LastfmUpdateDto {
  @ValidateNested()
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}