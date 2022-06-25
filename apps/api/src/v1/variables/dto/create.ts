import { CustomVar, CustomVarType } from '@tsuwari/prisma';
import { IsIn, IsNotEmpty, IsOptional, IsString, ValidateIf } from 'class-validator';

export class CreateVariableDto implements Partial<CustomVar> {
  @IsString()
  @IsNotEmpty()
  name: string;

  @IsString()
  @IsOptional()
  description?: string | null;

  @IsIn(Object.values(CustomVarType))
  type: CustomVarType;

  @ValidateIf((o: CreateVariableDto) => o.type === CustomVarType.SCRIPT)
  @IsString({ message: 'script should be string.' })
  @IsNotEmpty({ message: 'script should not be empty.' })
  evalValue?: string | null;

  @ValidateIf((o: CreateVariableDto) => o.type === CustomVarType.TEXT)
  @IsString()
  @IsNotEmpty({ message: 'text should be not empty.' })
  response?: string | null;
}