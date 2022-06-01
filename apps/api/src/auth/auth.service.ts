import { HttpException, HttpStatus, Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { ClientGrpc, ClientProxy } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { PrismaService } from '@tsuwari/prisma';
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

    const tokenData = {
      accessToken: tokens.accessToken,
      refreshToken: tokens.refreshToken!,
      obtainmentTimestamp: new Date(tokens.obtainmentTimestamp),
      expiresIn: tokens.expiresIn!,
    };

    const tokensQuery = {
      channel: {
        connectOrCreate: {
          where: {
            id: userId,
          },
          create: {
            isEnabled: true,
            botId: bot.id,
          },
        },
      },
      token: {
        upsert: {
          update: tokenData,
          create: tokenData,
        },
      },
    };

    const user = await this.prisma.user.upsert({
      where: { id: userId },
      update: tokensQuery,
      create: {
        id: userId,
        ...tokensQuery,
        token: {
          connectOrCreate: {
            where: { userId },
            create: tokenData,
          },
        },
      },
      include: {
        channel: true,
      },
    });

    if (username) {
      await this.botsMicroservice.joinOrLeave({
        action: user.channel?.isEnabled ? Bots.JoinOrLeaveRequest.Action.JOIN : Bots.JoinOrLeaveRequest.Action.PART,
        username,
      }).toPromise();
    }

    return user;
  }

  async getProfile(userId: string) {
    const [isTester, dashboards] = await Promise.all([
      this.prisma.tester.count({ where: { userId } }),
      this.prisma.dashboardAccess.findMany({
        where: { userId },
      }),
    ]);

    const neededUsers = await staticApi.users.getUsersByIds([
      ...dashboards.map(d => d.channelId),
      userId,
    ]);
    const user = neededUsers.find(u => u.id === userId);

    if (!user) throw new Error('User not found');

    return {
      ...getRawData(user),
      isTester: !!isTester,
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
