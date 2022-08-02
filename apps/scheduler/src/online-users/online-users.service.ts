import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { onlineUserSchema, RedisORMService } from '@tsuwari/redis';
import { TwitchApiService } from '@tsuwari/shared';
import _ from 'lodash';

@Injectable()
export class OnlineUsersService implements OnModuleInit {
  constructor(private readonly prisma: PrismaService, private readonly redisOm: RedisORMService, private readonly twitch: TwitchApiService) {

  }

  onModuleInit() {
    this.onlineUsers();
  }

  @Interval('onlineUsers', config.isDev ? 5000 : 1 * 60 * 1000)
  async onlineUsers() {
    const repository = this.redisOm.fetchRepository(onlineUserSchema);

    const channelsIds = await this.prisma.channel.findMany({
      where: {
        isEnabled: true,
        isBanned: false,
        isTwitchBanned: false,
      },
      select: {
        id: true,
      },
    });

    channelsIds.push({ id: '146712489' });

    const chunks = _.chunk(channelsIds.map(c => c.id), 100);

    for (const channelsIds of chunks) {
      const streams = await this.twitch.streams.getStreams({
        userId: channelsIds,
      });

      for (const stream of streams.data) {
        const users = await this.twitch.unsupported.getChatters(stream.userName);

        for (const user of users.allChatters) {
          repository.createAndSave({
            userName: user,
            userId: '',
            channelId: stream.userId,
          }, `${stream.userId}:${user}`).then(() => {
            repository.expire(`${stream.userId}:${user}`, 300);
          });
        }
      }
    }
  }
}
