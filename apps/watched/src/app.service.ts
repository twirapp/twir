import { Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { CachedStream } from '@tsuwari/shared';
import { ApiClient, HelixStream } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData, rawDataSymbol } from '@twurple/common';
import _ from 'lodash';


import { RedisService } from './redis.service.js';


const authProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  constructor(private readonly redis: RedisService, private readonly prisma: PrismaService) {

  }

  async handleIncrease(keys: string[]) {
    for (let index = 0; index < keys.length; index++) {
      const key = keys[index]!;

      const cachedStream = await this.redis.get(key);
      if (!cachedStream) continue;
      const stream = JSON.parse(cachedStream) as CachedStream;
      const chatters = await api.unsupported.getChatters(stream.user_login);

      const userNames = chatters.allChattersWithStatus.keys();
      const usersChunks = _.chunk([...userNames], 100);

      for (const chunk of usersChunks) {
        const users = await api.users.getUsersByNames(chunk);

        this.prisma.$transaction(
          users.map(user => {
            return this.prisma.userStats
              .upsert({
                where: {
                  userId_channelId: {
                    userId: user.id,
                    channelId: stream.user_id,
                  },
                },
                create: {
                  user: {
                    connectOrCreate: {
                      where: {
                        id: user.id,
                      },
                      create: {
                        id: user.id,
                      },
                    },
                  },
                  channel: {
                    connect: {
                      id: stream.user_id,
                    },
                  },
                },
                update: {
                  watched: {
                    increment: config.isDev ? 10000 : 5 * 60 * 1000,
                  },
                },
              });
          }),
        ).then(users => {
          for (const user of users) {
            const key = `usersStats:${stream.user_id}:${user.userId}`;
            this.redis.hset(key, 'watched', String(user.watched)).then(() => {
              this.redis.expire(key, 600);
            });
          }
        });
      }
    }
  }
}