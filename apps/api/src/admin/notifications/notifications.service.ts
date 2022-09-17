import { HttpException, Injectable } from '@nestjs/common';
import { TwitchApiService } from '@tsuwari/shared';
import { Notification } from '@tsuwari/typeorm/entities/Notification';

import { typeorm } from '../../index.js';
import { CreateNotificationDto } from './dto/create.js';

@Injectable()
export class NotificationsService {
  constructor(private readonly twitchApi: TwitchApiService) {}

  async create(data: CreateNotificationDto) {
    const query = {} as Partial<Notification>;

    if (data.userName) {
      const user = await this.twitchApi.users.getUserByName(data.userName.toLowerCase());

      if (!user) throw new HttpException(`Use ${data.userName} not found on twitch.`, 404);
      query.userId = user.id;
    }

    delete data.userName;

    return typeorm.getRepository(Notification).create({
      ...(data as Omit<CreateNotificationDto, 'userName'>),
      ...query,
      messages: data.messages,
    });
  }

  delete(id: string) {
    return typeorm.getRepository(Notification).delete({
      id,
    });
  }
}
