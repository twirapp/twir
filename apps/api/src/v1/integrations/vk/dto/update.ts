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
  userId: string;
}

export class VkUpdateDto {
  @ValidateNested()
  @ValidateIf((o: VkUpdateDto) => o.enabled)
  @Type(() => Data)
  data: Data;

  @IsBoolean()
  enabled: boolean;
}
