import { Injectable } from '@nestjs/common';
import { IsNull, Not } from '@tsuwari/typeorm';
import { Notification } from '@tsuwari/typeorm/entities/Notification';
import { UserViewedNotification } from '@tsuwari/typeorm/entities/UserViewedNotification';

import { typeorm } from '../../index.js';

@Injectable()
export class NotificationsService {
  getNotViewed(channelId: string) {
    return typeorm.getRepository(Notification).find({
      where: [
        { userId: IsNull() },
        { userId: channelId },
        { viewedNotifications: { userId: Not(channelId) } },
      ],
      order: {
        createdAt: 'DESC',
      },
      relations: {
        messages: true,
      },
    });
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
