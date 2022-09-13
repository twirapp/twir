import { Controller, Get, Res } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import { ClientProxyEvents, EventPattern, MessagePattern, type ClientProxyCommandPayload, type ClientProxyEventPayload } from '@tsuwari/shared';

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
  joinOrLeave(@Payload() data: ClientProxyEventPayload<'bots.joinOrLeave'>) {
    const bot = Bots.cache.get(data.botId);
    if (bot) {
      bot[data.action](data.username);
    }
  }

  @EventPattern('user.update')
  async userUpdate(@Payload() data: ClientProxyEventPayload<'user.update'>) {
    const channel = await prisma.channel.findFirst({
      where: { id: data.user_id },
    });

    if (channel?.isEnabled) {
      const bot = Bots.cache.get(channel.botId);
      if (bot) {
        bot.join(data.user_name);
      }
    }
  }

  // @MessagePattern('bots.deleteMessages')
  async deleteMessages(@Payload() data: ClientProxyCommandPayload<'bots.deleteMessages'>) {
    const channel = await prisma.channel.findFirst({
      where: { id: data.channelId },
    });

    if (!channel) return false;

    const bot = Bots.cache.get(channel?.botId);
    if (!bot) return false;

    for (const id of data.messageIds) {
      await bot.deleteMessage(data.channelName, id);
    }

    return true;
  }
}
