import { Controller, Get, Res } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import { EventPattern, type ClientProxyCommandPayload, type ClientProxyEventPayload } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';

import { Bots } from '../bots.js';
import { prometheus } from '../libs/prometheus.js';
import { typeorm } from '../libs/typeorm.js';

@Controller()
export class AppController {

  @Get('/metrics')
  async root(@Res() res: any) {
    res.contentType(prometheus.contentType);
    res.send(await prometheus.register.metrics());
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
    const channel = await typeorm.getRepository(Channel).findOneBy({
      id: data.user_id,
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
    const channel = await typeorm.getRepository(Channel).findOneBy({
      id: data.channelId,
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
