import { HttpException, Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { Channel, PrismaService, Token, User } from '@tsuwari/prisma';
import { AccessToken } from '@twurple/auth';
import { getRawData } from '@twurple/common';

import { staticApi } from '../twitchApi.js';

@Injectable()
export class AuthService implements OnModuleInit {
  private botsMicroservice: Bots.Main;

  constructor(private readonly prisma: PrismaService, @Inject('BOTS_MICROSERVICE') private readonly botsClient: ClientGrpc) { }

  onModuleInit() {
    this.botsMicroservice = this.botsClient.getService<Bots.Main>('Main');
  }

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

    console.log(user);
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
      await this.botsMicroservice.joinOrLeave({
        action: user.channel?.isEnabled ? Bots.JoinOrLeaveRequest.Action.JOIN : Bots.JoinOrLeaveRequest.Action.PART,
        username,
      }).toPromise();
    }

    return user;
  }

  async getProfile(userId: string) {
    const [dbUser, dashboards] = await Promise.all([
      this.prisma.user.findFirst({ where: { id: userId } }),
      this.prisma.dashboardAccess.findMany({
        where: { userId },
      }),
    ]);

    const neededUsers = await staticApi.users.getUsersByIds([
      ...dashboards.map(d => d.channelId),
      userId,
    ]);
    const user = neededUsers.find(u => u.id === userId);

    if (!user || !dbUser) throw new Error('User not found');

    return {
      ...getRawData(user),
      isTester: dbUser.isTester,
      dashboards: dashboards.map(d => {
        const twitchUser = neededUsers.find(u => u.id === d.channelId);
        if (!twitchUser) return;
        return {
          ...d,
          twitch: getRawData(twitchUser),
        };
      }),
    };
  }
}
