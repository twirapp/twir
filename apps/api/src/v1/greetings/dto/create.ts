import { IsBoolean, IsOptional, IsString } from 'class-validator';

export class GreetingCreateDto {
  @IsString()
  username: string;

  @IsString()
  text: string;

  @IsBoolean()
  @IsOptional()
  enabled?: boolean;
}
