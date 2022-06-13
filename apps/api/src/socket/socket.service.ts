import { Injectable } from '@nestjs/common';
import { Client, ClientProxy, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { EventParams, ClientToServerEvents, MyRefreshingProvider } from '@tsuwari/shared';
import { ApiClient } from '@twurple/api';

import { RedisService } from '../redis.service.js';

@Injectable()
export class SocketService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
  ) { }

  async isBotMod(opts: EventParams<ClientToServerEvents, 'isBotMod'>[0]) {
    const channel = await this.prisma.channel.findFirst({
      where: {
        id: opts.channelId,
      },
      include: {
        bot: true,
        user: {
          include: {
            token: true,
          },
        },
      },
    });

    if (!channel?.bot || !channel.user?.token) throw new Error('Missed bot or token on the channel');

    const authProvider = new MyRefreshingProvider({
      clientId: config.TWITCH_CLIENTID,
      clientSecret: config.TWITCH_CLIENTSECRET,
      token: channel.user.token,
      prisma: this.prisma,
    });

    const token = await authProvider.getAccessToken();

    if (!token?.scope.includes('moderation:read')) {
      return !!await this.redis.get(`isBotMod:${opts.channelName}`);
    }

    const api = new ApiClient({ authProvider });

    const mods = await api.moderation.getModerators(opts.channelId);
    const isMod = !!mods.data.find(m => m.userId === channel.bot.id);

    const redisKey = `isBotMod:${opts.channelName}`;
    if (isMod) {
      this.redis.set(redisKey, 'true');
    } else {
      this.redis.del(redisKey);
    }

    return { value: isMod };
  }

  async botJoinOrLeave(action: 'join' | 'part', opts: EventParams<ClientToServerEvents, 'botJoin'>[0]) {
    const channel = await this.prisma.channel.findFirst({
      where: {
        id: opts.channelId,
      },
    });

    if (!channel?.botId) return;

    await this.prisma.channel.update({
      where: { id: channel.id },
      data: {
        isEnabled: action === 'join',
      },
    });

    await this.nats.emit('bots.joinOrLeave', {
      action,
      botId: channel.botId,
      username: opts.channelName,
    }).toPromise();
  }
}
