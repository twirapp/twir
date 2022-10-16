import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { RedisService, TwitchApiService } from '@tsuwari/shared';
import { In, IsNull, Not } from '@tsuwari/typeorm';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import { UserOnline } from '@tsuwari/typeorm/entities/UserOnline';

import { typeorm } from '../index.js';

@Injectable()
export class OnlineUsersService implements OnModuleInit {
  constructor(private readonly twitch: TwitchApiService, private readonly redis: RedisService) {}

  onModuleInit() {
    this.onlineUsers();
  }

  @Interval('onlineUsers', config.isDev ? 5000 : 1 * 60 * 1000)
  async onlineUsers() {
    const streams = await typeorm.getRepository(ChannelStream).find();
    const usersRepository = typeorm.getRepository(UserOnline);

    await Promise.all(
      streams.map(async (stream) => {
        const users = await this.twitch.unsupported.getChatters(stream.userLogin);
        const chatters = new Set(users.allChatters);

        const current = await usersRepository.findBy({
          channelId: stream.userId,
          userName: Not(IsNull()),
        });
        const forDelete = current.filter((c) => !chatters.has(c.userName!));
        const forCreate = [...chatters.values()].filter(
          (c) => !current.some((cur) => cur.userName! === c),
        );

        await usersRepository.delete({
          channelId: stream.userId,
          userName: In(forDelete),
        });

        await usersRepository.save(
          forCreate.map((item) => ({
            channelId: stream.userId,
            userName: item,
          })),
        );
      }),
    );
  }
}
