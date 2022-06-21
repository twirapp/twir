import { Type } from 'class-transformer';
import { IsBoolean, IsString, ValidateNested } from 'class-validator';

class Data {
  @IsString()
  userId: string;
}

export class VkUpdateDto {
  @ValidateNested()
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}