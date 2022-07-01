import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { Channel, PrismaService, Token, User } from '@tsuwari/prisma';
import { ClientProxy, AuthUser } from '@tsuwari/shared';
import { HelixUser } from '@twurple/api';
import { AccessToken } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import chunk from 'lodash.chunk';

import { JwtPayload } from '../jwt/jwt.strategy.js';
import { staticApi } from '../twitchApi.js';

@Injectable()
export class AuthService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(private readonly prisma: PrismaService) { }

  async checkUser(tokens: AccessToken, userId: string, username?: string | null) {
    const bot = await this.prisma.bot.findFirst({
      where: {
        type: 'DEFAULT',
      },
    });

    if (!bot) {
      throw new Error('Bot not created, cannot create user.');
    }

    if (!tokens.refreshToken || !tokens.expiresIn) {
      throw new HttpException(`Something went wrong on gettings twitch tokens. Please, try again later.`, 500);
    }

    const tokenData = {
      accessToken: tokens.accessToken,
      refreshToken: tokens.refreshToken,
      obtainmentTimestamp: new Date(tokens.obtainmentTimestamp),
      expiresIn: tokens.expiresIn,
    };

    let user: (User & {
      channel: Channel | null;
      token: Token | null;
    }) | null = await this.prisma.user.findFirst({
      where: { id: userId },
      include: { channel: true, token: true },
    });

    if (user) {
      if (!user.channel) {
        user.channel = await this.prisma.channel.create({
          data: { id: user.id, botId: bot.id },
        });
      }

      if (user.tokenId) {
        await this.prisma.token.update({
          where: {
            id: user.tokenId,
          },
          data: tokenData,
        });
      } else {
        await this.prisma.user.update({
          where: { id: userId },
          data: {
            token: { create: tokenData },
          },
        });
      }
    } else {
      user = await this.prisma.user.create({
        data: {
          id: userId,
          channel: { create: { botId: bot.id } },
          token: { create: tokenData },
        },
        include: {
          channel: true,
          token: true,
        },
      });
    }

    if (username) {
      await this.nats.emit('bots.joinOrLeave', {
        action: user.channel?.isEnabled ? 'join' : 'part',
        username,
        botId: bot.id,
      }).toPromise();
    }

    return user;
  }

  async getProfile(userPayload: JwtPayload) {
    const [dbUser, dashboards] = await Promise.all([
      this.prisma.user.findFirst({ where: { id: userPayload.id } }),
      this.prisma.dashboardAccess.findMany({
        where: { userId: userPayload.id },
      }),
    ]);

    if (dbUser?.isBotAdmin) {
      const channels = await this.prisma.channel.findMany({
        where: {
          id: {
            notIn: [...dashboards.map(d => d.channelId), dbUser.id],
          },
        },
      });

      for (const channel of channels) {
        dashboards.push({
          id: channel.id,
          channelId: channel.id,
          userId: dbUser.id,
        });
      }
    }

    const chunks = chunk([...dashboards.map(d => d.channelId), userPayload.id], 100);
    const twitchUsers = await Promise.all(chunks.map((c) => staticApi.users.getUsersByIds(c))).then(v => v.flat());

    const user = twitchUsers.find(u => u.id === userPayload.id);

    if (!user || !dbUser) throw new HttpException('User not found', 404);

    const result: AuthUser = {
      ...getRawData(user),
      isTester: dbUser.isTester,
      dashboards: dashboards.map(d => {
        const twitchUser = twitchUsers.find(u => u.id === d.channelId);
        if (!twitchUser) return;
        return {
          ...d,
          twitch: getRawData(twitchUser),
        };
      }).filter(Boolean) as AuthUser['dashboards'],
    };


    if (dbUser.isBotAdmin) {
      result.isBotAdmin = dbUser.isBotAdmin;
    }

    return result;
  }
}
