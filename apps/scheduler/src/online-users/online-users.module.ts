import { Module } from '@nestjs/common';
import { TwitchApiService } from '@tsuwari/shared';

import { OnlineUsersService } from './online-users.service.js';

@Module({
  imports: [],
  providers: [TwitchApiService, OnlineUsersService],
})
export class OnlineUsersModule { }
