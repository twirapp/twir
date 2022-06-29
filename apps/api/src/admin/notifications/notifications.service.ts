import { HttpException, Injectable } from '@nestjs/common';
import { PrismaService, Notification } from '@tsuwari/prisma';
import { TwitchApiService } from '@tsuwari/shared';

import { CreateNotificationDto } from './dto/create.js';

@Injectable()
export class NotificationsService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly twitchApi: TwitchApiService,
  ) { }

  async create(data: CreateNotificationDto) {
    const query = {} as Partial<Notification>;

    if (data.userName) {
      const user = await this.twitchApi.users.getUserByName(data.userName.toLowerCase());

      if (!user) throw new HttpException(`Use ${data.userName} not found on twitch.`, 404);
      query.userId = user.id;
    }

    delete data.userName;

    return this.prisma.notification.create({
      data: {
        ...data as Omit<CreateNotificationDto, 'userName'>,
        ...query,
      },
    });
  }

  delete(id: string) {
    return this.prisma.notification.delete({
      where: { id },
    });
  }
}
