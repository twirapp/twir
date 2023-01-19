import { Injectable, OnModuleInit } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { In } from '@tsuwari/typeorm';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import { IgnoredUser } from '@tsuwari/typeorm/entities/IgnoredUser';
import { User } from '@tsuwari/typeorm/entities/User';
import { UserOnline } from '@tsuwari/typeorm/entities/UserOnline';
import { ApiClient } from '@twurple/api';
import { StaticAuthProvider } from '@twurple/auth';
import _ from 'lodash';

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
        const userToken = await tokensGrpcClient.requestUserToken({
          userId: stream.userId,
        })
          .then((t) => {
            if (!t.scopes.includes('moderator:read:chatters')) {
              return null;
            } else return t;
          })
          .catch(() => null);

        const allChatters: Array<{ channelId: string, userId: string, userName?: string }> = [];

        if (userToken) {
          const apiClient = new ApiClient({
            authProvider: new StaticAuthProvider(config.TWITCH_CLIENTID, userToken.accessToken, userToken.scopes),
          });

          const getChatters = async(after?: string) => {
            const users = await apiClient.chat.getChatters(stream.userId, stream.userId, { after });
            return {
              data: users.data.map(u => ({
                userId: u.userId,
                userName: u.userDisplayName,
                channelId: stream.userId,
              })),
              cursor: users.cursor,
            };
          };

          let cursor: string | undefined;

          // eslint-disable-next-line no-constant-condition
          while (true) {
            const chatters = await getChatters(cursor);
            allChatters.push(...chatters.data);
            if (!chatters.cursor) {
              break;
            } else {
              cursor = chatters.cursor;
            }
          }
        } else {
          const chatters = await apiClient.unsupported.getChatters(stream.userLogin);
          const chunks = _.chunk(chatters.allChatters, 100);

          const users = await Promise.all(chunks.map(async (c) => {
            const users = await apiClient.users.getUsersByNames(c);
            return users;
          }));

          const mappedUsers = users
            .flat()
            .filter(c => c.id)
            .map(u => ({
              userId: u.id,
              userName: u.displayName,
              channelId: stream.userId,
            }));
          allChatters.push(...mappedUsers);
        }

        const ignoredOnlineChatters = await typeorm
          .getRepository(IgnoredUser)
          .findBy({
            id: In(allChatters.map(c => c.userId)),
          });

        const filteredAllChatters = allChatters.filter(c => !ignoredOnlineChatters.some(s => s.id === c.userId));
        const current = await usersRepository.findBy({
          channelId: stream.userId,
        });

        const forDelete = current.filter((curr) => !filteredAllChatters.some(c => c.userId === curr.userId));
        const forCreate = filteredAllChatters.filter(
          (c) => !current.some((cur) => cur.userId === c.userId),
        );

        const forCreateChunks = _.chunk(forCreate, 1000);

        await typeorm.transaction(async (manager) => {
          await manager.remove(forDelete);
          for (const chunk of forCreateChunks) {
            const repository = typeorm.getRepository(User);
            const existedUsers = await repository.findBy({
              id: In(chunk.map(u => u.userId)),
            });

            const usersForCreate = chunk.filter(u => !existedUsers.some(e => e.id === u.userId));
            await manager.save(usersForCreate.map(u => repository.create({ id: u.userId })));
            await manager.save(chunk.map(u => usersRepository.create(u)));
          }
        });
      }),
    );
  }
}
