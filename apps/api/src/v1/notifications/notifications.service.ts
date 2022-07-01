import { Injectable } from '@nestjs/common';
import { LangCode, PrismaService } from '@tsuwari/prisma';

@Injectable()
export class NotificationsService {
  constructor(private readonly prisma: PrismaService) {

  }

  getNotViewed(channelId: string) {
    return this.prisma.notification.findMany({
      where: {
        OR: [{ userId: null }, { userId: channelId }],
        viewedNotifications: {
          none: {
            userId: channelId,
          },
        },
      },
      include: {
        messages: true,
      },
      orderBy: {
        createdAt: 'desc',
      },
    });
  }

  getViewed(channelId: string) {
    return this.prisma.viewedNotification.findMany({
      where: {
        userId: channelId,
      },
      include: {
        notification: {
          include: {
            messages: true,
          },
        },
      },
      orderBy: {
        createdAt: 'desc',
      },
    });
  }

  markAsRead(channelId: string, notificationId: string) {
    return this.prisma.viewedNotification.create({
      data: {
        userId: channelId,
        notificationId,
      },
    });
  }
}
