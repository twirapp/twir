import { LangCode } from '@tsuwari/prisma';
import { Type } from 'class-transformer';
import { IsArray, IsIn, IsNotEmpty, IsOptional, IsString, ValidateNested } from 'class-validator';

export class CreateNotificationMessage {
  @IsString()
  @IsOptional()
  @IsNotEmpty()
  title?: string;

  @IsString()
  @IsNotEmpty()
  text: string;

  @IsIn(Object.values(LangCode))
  langCode: LangCode;
}

export class CreateNotificationDto {
  @IsString()
  @IsOptional()
  @IsNotEmpty()
  imageSrc?: string;

  @IsString()
  @IsOptional()
  @IsNotEmpty()
  userName?: string;

  @IsArray()
  @ValidateNested()
  @Type(() => CreateNotificationMessage)
  messages: Array<CreateNotificationMessage>;
}