import { IsNotEmpty, IsOptional, IsString } from 'class-validator';

export class CreateNotificationDto {
  @IsString()
  @IsOptional()
  @IsNotEmpty()
  title?: string;

  @IsString()
  @IsNotEmpty()
  text: string;

  @IsString()
  @IsOptional()
  @IsNotEmpty()
  imageSrc?: string;

  @IsString()
  @IsOptional()
  @IsNotEmpty()
  userName?: string;
}