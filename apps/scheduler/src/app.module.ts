import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { RedisModule, RedisService, TwitchApiService } from '@tsuwari/shared';

import { DeleteOldMessagesModule } from './delete-old-messages/deleteoldmessages.module.js';
import { EmotesModule } from './emotes/emotes.module.js';
import { MicroservicesModule } from './microservices/microservices.module.js';
import { OnlineUsersModule } from './online-users/online-users.module.js';
import { WatchedModule } from './watched/watched.module.js';

@Module({
  imports: [
    RedisModule,
    ScheduleModule.forRoot(),
    MicroservicesModule,
    OnlineUsersModule,
    DeleteOldMessagesModule,
    WatchedModule,
    EmotesModule,
  ],
  providers: [TwitchApiService, RedisService],
})
export class AppModule {}
