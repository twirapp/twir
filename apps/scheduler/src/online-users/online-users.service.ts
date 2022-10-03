import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { RedisService, TwitchApiService } from '@tsuwari/shared';

@Injectable()
export class OnlineUsersService implements OnModuleInit {
  constructor(private readonly twitch: TwitchApiService, private readonly redis: RedisService) {}

  onModuleInit() {
    this.onlineUsers();
  }

  @Interval('onlineUsers', config.isDev ? 5000 : 1 * 60 * 1000)
  async onlineUsers() {
    const streamsKeys = await this.redis.keys('streams:*');
    if (!streamsKeys.length) return;

    const streams = await this.redis.mget(streamsKeys);

    for (const streamData of streams.filter((v) => v != null)) {
      const stream = JSON.parse(streamData!);
      const users = await this.twitch.unsupported.getChatters(stream.user_login);

      for (const user of users.allChatters) {
        this.redis.set(
          `onlineUsers:${stream.user_id}:${user}`,
          JSON.stringify({
            username: user,
            userId: '',
            channelId: stream.user_id,
          }),
          'EX',
          300,
        );
      }
    }
  }
}
