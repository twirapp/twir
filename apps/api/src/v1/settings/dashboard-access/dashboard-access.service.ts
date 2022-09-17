import { HttpException, Injectable } from '@nestjs/common';
import { Not } from '@tsuwari/typeorm';
import { DashboardAccess } from '@tsuwari/typeorm/entities/DashboardAccess';
import { User } from '@tsuwari/typeorm/entities/User';
import { getRawData } from '@twurple/common';

import { typeorm } from '../../../index.js';
import { staticApi } from '../../../twitchApi.js';

@Injectable()
export class DashboardAccessService {
  async #ensureUser(id: string) {
    const repository = typeorm.getRepository(User);

    if (!(await repository.findOneBy({ id }))) {
      await repository.save({
        id,
      });
    }
  }

  async getMembers(channelId: string) {
    const members = await typeorm.getRepository(DashboardAccess).find({
      where: {
        channelId,
        userId: Not(channelId),
      },
    });

    const twitchUsers = await staticApi.users.getUsersByIds(members.map((m) => m.userId));

    return twitchUsers.map((u) => getRawData(u));
  }

  async addMember(channelId: string, username: string) {
    const twitchUser = await staticApi.users.getUserByName(username);

    if (!twitchUser)
      throw new HttpException(`Member with username ${username} not found on twitch.`, 404);
    if (channelId === twitchUser.id)
      throw new HttpException(`You cannot add youself as member`, 400);

    const repository = typeorm.getRepository(DashboardAccess);

    const isExists = await repository.findOneBy({
      channelId,
      userId: twitchUser.id,
    });

    if (isExists) throw new HttpException(`Member with name ${username} already exists.`, 400);

    await this.#ensureUser(twitchUser.id);
    await repository.save({
      channelId,
      userId: twitchUser.id,
    });
  }

  async deleteMember(channelId: string, userId: string) {
    const repository = typeorm.getRepository(DashboardAccess);
    const entity = await repository.findOneBy({
      channelId,
      userId: userId,
    });

    if (!entity) throw new HttpException(`Member with id not member of dashboard.`, 404);

    await repository.delete({
      id: entity.id,
    });
  }
}
