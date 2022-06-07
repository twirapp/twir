import { CustomVar, CustomVarType } from '@tsuwari/prisma';
import { IsIn, IsNotEmpty, IsOptional, IsString, ValidateIf } from 'class-validator';

export class CreateVariableDto implements Partial<CustomVar> {
  @IsString()
  name: string;

  @IsString()
  @IsOptional()
  description?: string | null;

  @IsIn(Object.values(CustomVarType))
  type: CustomVarType;

  @ValidateIf((o: CreateVariableDto) => o.type === CustomVarType.SCRIPT)
  @IsString()
  @IsNotEmpty()
  evalValue?: string | null;

  @ValidateIf((o: CreateVariableDto) => o.type === CustomVarType.TEXT)
  @IsString()
  @IsNotEmpty()
  response?: string | null;
}