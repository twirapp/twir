import {
  IsBoolean,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsString,
  Max,
  MaxLength,
  Min,
} from 'class-validator';

export class CreateKeywordDto {
  @IsString()
  @IsNotEmpty()
  @MaxLength(100)
  text: string;

  @IsString()
  @IsNotEmpty()
  @MaxLength(400)
  response: string;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;

  @IsNumber()
  @Min(5)
  @Max(84000)
  @IsOptional()
  cooldown?: number;

  @IsBoolean()
  @IsOptional()
  isReply?: boolean;
}
