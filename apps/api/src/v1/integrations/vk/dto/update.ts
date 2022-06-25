import { Type } from 'class-transformer';
import { IsBoolean, IsNotEmpty, IsString, MaxLength, ValidateNested } from 'class-validator';

class Data {
  @IsString()
  @IsNotEmpty()
  @MaxLength(100)
  userId: string;
}

export class VkUpdateDto {
  @ValidateNested()
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}