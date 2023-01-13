import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { IsNull, Not } from '@tsuwari/typeorm';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import { UserOnline } from '@tsuwari/typeorm/entities/UserOnline';
import { ApiClient } from '@twurple/api';
import { StaticAuthProvider } from '@twurple/auth';

import { typeorm } from '../index.js';
import { tokensGrpcClient } from '../libs/tokens.grpc.js';

@Injectable()
export class OnlineUsersService implements OnModuleInit {
  onModuleInit() {
    this.onlineUsers();
  }

  @Interval('onlineUsers', config.isDev ? 5000 : 1 * 60 * 1000)
  async onlineUsers() {
    const streams = await typeorm.getRepository(ChannelStream).find();
    const usersRepository = typeorm.getRepository(UserOnline);

    const appToken = await tokensGrpcClient.requestAppToken({});

    const apiClient = new ApiClient({
      authProvider: new StaticAuthProvider(config.TWITCH_CLIENTID, appToken.accessToken),
    });

    await Promise.all(
      streams.map(async (stream) => {
        const users = await apiClient.unsupported.getChatters(stream.userLogin);
        const chatters = new Set(users.allChatters);

        const current = await usersRepository.findBy({
          channelId: stream.userId,
          userName: Not(IsNull()),
        });
        const forDelete = current.filter((c) => !chatters.has(c.userName!));
        const forCreate = [...chatters.values()].filter(
          (c) => !current.some((cur) => cur.userName! === c),
        );

        await usersRepository.remove(forDelete);
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
