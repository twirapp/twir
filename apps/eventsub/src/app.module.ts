import { Module } from '@nestjs/common';
import { TwitchApiService } from '@tsuwari/shared';

import { EventSubModule } from './eventsub/eventsub.module.js';
import { HandlerModule } from './handler/handler.module.js';

@Module({
  imports: [EventSubModule.register(), HandlerModule],
  controllers: [],
  providers: [TwitchApiService],
})
export class AppModule {}
