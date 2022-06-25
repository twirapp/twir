import { IsBoolean, IsNotEmpty, IsOptional, IsString } from 'class-validator';

export class GreetingCreateDto {
  @IsString()
  @IsNotEmpty()
  username: string;

  @IsString()
  @IsNotEmpty()
  text: string;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;
}
