import { Module } from '@nestjs/common';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisModule, RedisService, TwitchApiService } from '@tsuwari/shared';

import { EventSubModule } from './eventsub/eventsub.module.js';
import { HandlerModule } from './handler/handler.module.js';

@Module({
  imports: [EventSubModule.register(), PrismaModule, HandlerModule, RedisModule],
  controllers: [],
  providers: [TwitchApiService, RedisService],
})
export class AppModule { }
