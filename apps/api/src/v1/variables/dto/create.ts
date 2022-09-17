import { ChannelCustomvar, CustomVarType } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { IsIn, IsNotEmpty, IsOptional, IsString, MaxLength, ValidateIf } from 'class-validator';

export class CreateVariableDto implements Partial<ChannelCustomvar> {
  @IsString()
  @IsNotEmpty()
  @MaxLength(20)
  name: string;

  @IsString()
  @IsOptional()
  @MaxLength(400)
  description?: string | null;

  @IsIn(Object.values(CustomVarType))
  type: CustomVarType;

  @ValidateIf((o: CreateVariableDto) => o.type === CustomVarType.SCRIPT)
  @IsString({ message: 'script should be string.' })
  @IsNotEmpty({ message: 'script should not be empty.' })
  @MaxLength(10000)
  evalValue?: string | null;

  @ValidateIf((o: CreateVariableDto) => o.type === CustomVarType.TEXT)
  @IsString()
  @IsNotEmpty({ message: 'text should be not empty.' })
  @MaxLength(500)
  response?: string | null;
}
