import { Controller, Get, Res } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';

import { Bots } from '../bots.js';
import { prisma } from '../libs/prisma.js';
import { prometheus } from '../libs/prometheus.js';

@Controller()
export class AppController {
  @Get('/metrics')
  async root(@Res() res: any) {
    res.contentType(prometheus.contentType);
    res.send(await prometheus.register.metrics() + await prisma.$metrics.prometheus());
  }

  @EventPattern('bots.joinOrLeave')
  joinOrLeave(@Payload() data: { action: 'join' | 'part', botId: string, username: string }) {
    const bot = Bots.cache.get(data.botId);
    if (bot) {
      bot[data.action](data.username);
    }
  }
}
