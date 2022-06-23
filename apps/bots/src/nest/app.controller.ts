import { Controller } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';

import { Bots } from '../bots.js';

@Controller()
export class AppController {
  @EventPattern('bots.joinOrLeave')
  joinOrLeave(@Payload() data: { action: 'join' | 'part', botId: string, username: string }) {
    const bot = Bots.cache.get(data.botId);
    if (bot) {
      bot[data.action](data.username);
    }
  }
}
