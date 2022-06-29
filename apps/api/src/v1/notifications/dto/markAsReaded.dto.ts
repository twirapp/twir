import { IsNotEmpty, IsString, MaxLength } from 'class-validator';

export class MarkAsReadedDto {
  @IsString()
  @IsNotEmpty()
  @MaxLength(100)
  notificationId: string;
}