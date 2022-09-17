import { Injectable } from '@nestjs/common';
import { ArrayContains, In, IsNull, Not } from '@tsuwari/typeorm';
import { Notification } from '@tsuwari/typeorm/entities/Notification';
import { UserViewedNotification } from '@tsuwari/typeorm/entities/UserViewedNotification';

import { typeorm } from '../../index.js';

@Injectable()
export class NotificationsService {
  async getNotViewed(channelId: string) {
    const viewed = await typeorm.getRepository(UserViewedNotification).find({
      where: {
        userId: channelId,
      },
    });

    const mappedViewed = viewed.map((v) => v.notificationId);

    const notifications = await typeorm.getRepository(Notification).find({
      where: [
        { id: Not(In(mappedViewed)), userId: IsNull() },
        { id: Not(In(mappedViewed)), userId: channelId },
      ],
      order: {
        createdAt: 'DESC',
      },
      relations: {
        messages: true,
      },
    });

    return notifications;
  }

  getViewed(channelId: string) {
    return typeorm.getRepository(UserViewedNotification).find({
      where: { userId: channelId },
      relations: {
        notification: {
          messages: true,
        },
      },
      order: {
        createdAt: 'DESC',
      },
    });
  }

  markAsRead(channelId: string, notificationId: string) {
    return typeorm.getRepository(UserViewedNotification).save({
      userId: channelId,
      notificationId,
    });
  }
}
