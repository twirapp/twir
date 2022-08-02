import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { onlineUserSchema, RedisORMService, streamSchema } from '@tsuwari/redis';
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
    const streamsRepisitory = this.redisOm.fetchRepository(streamSchema);

    const streams = await streamsRepisitory.search().all();

    for (const streamData of streams) {
      const stream = streamData.toRedisJson();
      const users = await this.twitch.unsupported.getChatters(stream.user_login);

      for (const user of users.allChatters) {
        repository.createAndSave({
          userName: user,
          userId: '',
          channelId: stream.user_id,
        }, `${stream.user_id}:${user}`).then(() => {
          repository.expire(`${stream.user_id}:${user}`, 300);
        });
      }
    }
  }
}
