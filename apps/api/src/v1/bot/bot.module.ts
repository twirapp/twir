import { Module } from '@nestjs/common';
import { TwitchApiService } from '@tsuwari/shared';

import { BotController } from './bot.controller.js';
import { BotService } from './bot.service.js';

@Module({
  controllers: [BotController],
  providers: [TwitchApiService, BotService],
})
export class BotModule {}
