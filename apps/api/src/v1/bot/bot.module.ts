import { Module } from '@nestjs/common';

import { BotController } from './bot.controller.js';
import { BotService } from './bot.service.js';

@Module({
  controllers: [BotController],
  providers: [BotService],
})
export class BotModule {}
