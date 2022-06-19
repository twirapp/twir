import { HttpException, Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { getRawData } from '@twurple/common';

import { staticApi } from '../../../twitchApi.js';

@Injectable()
export class DashboardAccessService {
  constructor(private readonly prisma: PrismaService) { }

  async #ensureUser(id: string) {
    if (!await this.prisma.user.findFirst({ where: { id } })) {
      await this.prisma.user.create({
        data: { id },
      });
    }
  }

  async getMembers(channelId: string) {
    const members = await this.prisma.dashboardAccess.findMany({
      where: {
        channelId,
        AND: {
          userId: {
            not: channelId,
          },
        },
      },
    });

    const twitchUsers = await staticApi.users.getUsersByIds(members.map(m => m.userId));

    return twitchUsers.map(u => getRawData(u));
  }

  async addMember(channelId: string, username: string) {
    const twitchUser = await staticApi.users.getUserByName(username);

    if (!twitchUser) throw new HttpException(`Member with username ${username} not found on twitch.`, 404);

    const isExists = await this.prisma.dashboardAccess.findFirst({
      where: {
        channelId,
        userId: twitchUser.id,
      },
    });

    if (isExists) throw new HttpException(`Member with name ${username} already exists.`, 400);

    await this.#ensureUser(twitchUser.id);
    await this.prisma.dashboardAccess.create({
      data: {
        channelId,
        userId: twitchUser.id,
      },
    });
  }

  async deleteMember(channelId: string, userId: string) {
    const isExists = await this.prisma.dashboardAccess.findFirst({
      where: {
        channelId,
        userId: userId,
      },
    });

    if (!isExists) throw new HttpException(`Member with id not member of dashboard.`, 404);

    await this.prisma.dashboardAccess.delete({
      where: {
        id: isExists.id,
      },
    });
  }
}
