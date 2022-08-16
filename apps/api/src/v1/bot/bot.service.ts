import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { ClientProxy, MyRefreshingProvider, RedisService, TwitchApiService } from '@tsuwari/shared';
import { ApiClient } from '@twurple/api';

@Injectable()
export class BotService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly twitchApi: TwitchApiService,
  ) {}

  async isBotMod(channelId: string) {
    const channel = await this.prisma.channel.findFirst({
      where: {
        id: channelId,
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

    if (!channel?.bot || !channel.user?.token)
      throw new Error('Missed bot or token on the channel');

    const authProvider = new MyRefreshingProvider({
      clientId: config.TWITCH_CLIENTID,
      clientSecret: config.TWITCH_CLIENTSECRET,
      token: channel.user.token,
      prisma: this.prisma,
    });

    const token = await authProvider.getAccessToken();

    if (!token?.scope.includes('moderation:read')) {
      return !!(await this.redis.get(`isBotMod:${channelId}`));
    }

    const api = new ApiClient({ authProvider });

    const mods = await api.moderation.getModerators(channelId);
    const isMod = !!mods.data.find((m) => m.userId === channel.bot.id);

    const redisKey = `isBotMod:${channelId}`;
    if (isMod) {
      this.redis.set(redisKey, 'true');
    } else {
      this.redis.del(redisKey);
    }

    return isMod;
  }

  async botJoinOrLeave(action: 'join' | 'part', channelId: string) {
    const [channel, user] = await Promise.all([
      this.prisma.channel.findFirst({
        where: {
          id: channelId,
        },
      }),
      this.twitchApi.users.getUserByIdWithCache(channelId),
    ]);

    if (!user || !channel) throw new HttpException(`User not found`, 404);
    if (!channel.botId) {
      const defaultBot = await this.prisma.bot.findFirst({
        where: { type: 'DEFAULT' },
      });
      if (defaultBot) {
        await this.prisma.channel.update({
          where: {
            id: channel.id,
          },
          data: {
            botId: defaultBot.id,
          },
        });
      }
    }

    await Promise.all([
      this.prisma.channel.update({
        where: { id: channel.id },
        data: {
          isEnabled: action === 'join',
        },
      }),
      this.nats
        .emit('bots.joinOrLeave', {
          action,
          botId: channel.botId,
          username: user.login,
        })
        .toPromise(),
    ]);
  }
}
