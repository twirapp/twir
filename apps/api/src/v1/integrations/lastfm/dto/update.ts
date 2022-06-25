
import { Type } from 'class-transformer';
import { IsBoolean, IsNotEmpty, IsString, ValidateNested } from 'class-validator';

class Data {
  @IsString()
  @IsNotEmpty()
  username: string;
}

export class LastfmUpdateDto {
  @ValidateNested()
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}