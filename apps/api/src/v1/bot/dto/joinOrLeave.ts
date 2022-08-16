import { IsIn, IsString } from 'class-validator';

export class JoinOrLeaveDto {
  @IsString()
  @IsIn(['join', 'part'])
  action: 'join' | 'part';
}
