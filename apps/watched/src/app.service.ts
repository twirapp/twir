import { Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaService } from '@tsuwari/prisma';
import { CachedStream } from '@tsuwari/shared';
import { ApiClient, HelixUserData } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import _ from 'lodash';


import { RedisService } from './redis.service.js';


const authProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
const api = new ApiClient({ authProvider });

@Injectable()
export class AppService {
  private readonly incrNumber = config.isDev ? 10000 : 5 * 60 * 1000;

  constructor(private readonly redis: RedisService, private readonly prisma: PrismaService) {

  }

  async #getCachedTwitchUsers() {
    const usersKeys = await this.redis.keys(`twitch:users:*`);
    const users: Array<HelixUserData> = [];

    for (const key of usersKeys) {
      const user = await this.redis.hgetall(`twitch:users:${key.split(':')[2]}`);
      if (Object.keys(user).length) users.push(user as unknown as HelixUserData);
    }

    return users;
  }

  async #cacheTwitchUsers(users: HelixUserData[]) {
    for (const user of users) {
      this.redis.hmset(`twitch:users:${user.id}`, user);
    }
  }

  async handleIncrease(keys: string[]) {
    for (let index = 0; index < keys.length; index++) {
      const key = keys[index]!;

      const cachedStream = await this.redis.get(key);
      if (!cachedStream) continue;
      const stream = JSON.parse(cachedStream) as CachedStream;
      const chatters = await api.unsupported.getChatters(stream.user_login);

      const cachedUsers = await this.#getCachedTwitchUsers();

      const userNamesForGetFromTwitch = [...chatters.allChattersWithStatus.keys()].filter(name => cachedUsers.some(u => u.login !== name));
      const usersChunks = _.chunk(userNamesForGetFromTwitch, 100);

      const twitchUsers = await Promise.all(usersChunks.map(c => api.users.getUsersByNames(c)))
        .then(v => v.flat())
        .then(v => {
          return v.map(u => getRawData(u));
        });
      this.#cacheTwitchUsers(twitchUsers);

      const usersForIncrease = [...cachedUsers, ...twitchUsers];

      for (let index = 0; index < usersForIncrease.length; index++) {
        const user = usersForIncrease[index];
        if (!user) continue;

        const key = `usersStats:${stream.user_id}:${user.id}`;
        const cachedUser = await this.redis.hgetall(key);

        if (Object.keys(cachedUser).length) {
          this.redis.hset(key, 'watched', Number(Number(cachedUser.watched) + this.incrNumber)).then(() => {
            this.redis.expire(key, 1200);
          });

          this.prisma.userStats.update({
            where: {
              userId_channelId: {
                userId: user.id,
                channelId: stream.user_id,
              },
            },
            data: {
              watched: { increment: this.incrNumber },
            },
          });
        } else {
          this.prisma.userStats.upsert({
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
            select: {
              userId: true,
              watched: true,
            },
          }).then(u => {
            this.redis.hset(key, 'watched', String(u.watched)).then(() => {
              this.redis.expire(key, 1200);
            });
          });
        }
      }

      /* await Promise.all(_.chunk(twitchUsers, 1000).map(c => {
        return this.prisma.$transaction(c.map(user => {
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
        }));
      })).then(users => {
        for (const user of users.flat()) {
          const key = `usersStats:${stream.user_id}:${user.userId}`;
          this.redis.hset(key, 'watched', String(user.watched)).then(() => {
            this.redis.expire(key, 600);
          });
        }

        console.log('done');
      }); */
    }
  }
}