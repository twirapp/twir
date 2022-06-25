
import { Type } from 'class-transformer';
import { IsBoolean, IsNotEmpty, IsString, MaxLength, ValidateNested } from 'class-validator';

class Data {
  @IsString()
  @IsNotEmpty()
  @MaxLength(100)
  username: string;
}

export class LastfmUpdateDto {
  @ValidateNested()
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}