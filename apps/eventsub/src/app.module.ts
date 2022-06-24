import { Module } from '@nestjs/common';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisService, TwitchApiService } from '@tsuwari/shared';

import { EventSubModule } from './eventsub/eventsub.module.js';
import { HandlerModule } from './handler/handler.module.js';

@Module({
  imports: [EventSubModule.register(), PrismaModule, HandlerModule],
  controllers: [],
  providers: [TwitchApiService, RedisService],
})
export class AppModule { }
