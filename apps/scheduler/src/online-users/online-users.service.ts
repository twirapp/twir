import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { onlineUserSchema, RedisORMService } from '@tsuwari/redis';
import { RedisService, TwitchApiService } from '@tsuwari/shared';

@Injectable()
export class OnlineUsersService implements OnModuleInit {
  constructor(
    private readonly redisOm: RedisORMService,
    private readonly twitch: TwitchApiService,
    private readonly redis: RedisService,
  ) {}

  onModuleInit() {
    this.onlineUsers();
  }

  @Interval('onlineUsers', config.isDev ? 5000 : 1 * 60 * 1000)
  async onlineUsers() {
    const repository = this.redisOm.fetchRepository(onlineUserSchema);

    const streams = await this.redis.mget('streams:*');

    for (const streamData of streams.filter((v) => v != null)) {
      const stream = JSON.parse(streamData!);
      const users = await this.twitch.unsupported.getChatters(stream.user_login);

      for (const user of users.allChatters) {
        repository
          .createAndSave(
            {
              userName: user,
              userId: '',
              channelId: stream.user_id,
            },
            `${stream.user_id}:${user}`,
          )
          .then(() => {
            repository.expire(`${stream.user_id}:${user}`, 300);
          });
      }
    }
  }
}
