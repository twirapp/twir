import { CommandPermission, CooldownType } from '@tsuwari/prisma';
import { Type } from 'class-transformer';
import { ArrayNotEmpty, IsArray, IsBoolean, IsIn, IsNotEmpty, IsNumber, IsOptional, IsString, Max, MaxLength, Min, ValidateIf, ValidateNested } from 'class-validator';

export class UpdateOrCreateCommandDto {
  id?: string;
  default?: boolean;

  @IsString()
  @IsNotEmpty()
  @MaxLength(50)
  name: string;

  @ValidateIf((o: UpdateOrCreateCommandDto) => o.default === false)
  @IsNumber()
  @IsOptional()
  @Min(5)
  @Max(86400, { always: true })
  cooldown?: number;

  @IsIn(Object.keys(CooldownType))
  @IsOptional()
  cooldownType?: CooldownType;

  @IsString()
  @IsOptional()
  @MaxLength(400, { each: true })
  description?: string | null;

  @IsIn(Object.keys(CommandPermission))
  permission: CommandPermission;

  @IsArray()
  @IsOptional()
  @IsNotEmpty({ each: true })
  @MaxLength(50, { each: true })
  aliases?: string[];

  @IsBoolean()
  @IsOptional()
  visible?: boolean;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;

  @ValidateIf((o: UpdateOrCreateCommandDto) => o.default === false)
  @IsArray()
  @ValidateNested()
  @Type(() => ResponseValidation)
  @ArrayNotEmpty()
  responses: Array<ResponseValidation>;
}


class ResponseValidation {
  id?: string;

  @MaxLength(400, { each: true })
  @IsNotEmpty({ each: true })
  text: string | null;
}
