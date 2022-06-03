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

    const request = await Promise.all(usersKeys.map(k => this.redis.hgetall(`twitch:users:${k.split(':')[2]}`)));
    return request.filter(d => Object.keys(d).length > 0) as unknown as Array<HelixUserData>;
  }

  #cacheTwitchUsers(users: HelixUserData[]) {
    for (const user of users) {
      const key = `twitch:users:${user.id}`;
      this.redis.hmset(key, user).then(() => {
        this.redis.expire(key, 5 * 3600);
      });
    }
  }

  async handleIncrease(keys: string[]) {
    const start = performance.now();
    const cachedUsers = await this.#getCachedTwitchUsers();

    await Promise.all(keys.map(async (key) => {
      const cachedStream = await this.redis.get(key);
      if (!cachedStream) return;
      const stream = JSON.parse(cachedStream) as CachedStream;
      const chatters = await api.unsupported.getChatters(stream.user_login);

      const cachedUsersNames = cachedUsers.map(u => u.login);

      const userNamesForGetFromTwitch = [...chatters.allChattersWithStatus.keys()].filter(name => !cachedUsersNames.includes(name));
      const usersChunks = _.chunk(userNamesForGetFromTwitch, 100);

      const twitchUsers = await Promise.all(usersChunks.map(c => api.users.getUsersByNames(c)))
        .then(v => v.flat())
        .then(v => {
          return v.map(u => getRawData(u));
        });
      this.#cacheTwitchUsers(twitchUsers);

      const usersForIncrease = [...cachedUsers, ...twitchUsers];
      const usersForUpdate: string[] = [];
      const usersForUpsert: string[] = [];

      const cachedUsersStats = await Promise.all(usersForIncrease.map(async (u) => {
        return {
          id: u.id,
          data: await this.redis.hgetall(`usersStats:${stream.user_id}:${u.id}`),
        };
      }));

      cachedUsersStats.forEach(user => {
        const userKey = `usersStats:${stream.user_id}:${user.id}`;
        if (Object.keys(user.data).length) {
          this.redis.hset(`usersStats:${stream.user_id}:${user.id}`, 'watched', Number(Number(user.data.watched) + this.incrNumber)).then(() => {
            this.redis.expire(userKey, 1200);
          });

          usersForUpdate.push(user.id);
        } else {
          usersForUpsert.push(user.id);
        }
      });

      /* this.prisma.userStats.updateMany({
        where: {
          channelId: stream.user_id,
          userId: { in: usersForUpdate },
        },
        data: {
          watched: { increment: this.incrNumber },
        },
      });

      this.prisma.$transaction(usersForUpsert.map(u => {
        return this.prisma.userStats.upsert({
          where: {
            userId_channelId: {
              userId: u,
              channelId: stream.user_id,
            },
          },
          create: {
            user: {
              connectOrCreate: {
                where: {
                  id: u,
                },
                create: {
                  id: u,
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
              increment: this.incrNumber,
            },
          },
          select: {
            userId: true,
            watched: true,
          },
        });
      })).then(users => {
        for (const user of users) {
          this.redis.hset(`usersStats:${stream.user_id}:${user.userId}`, 'watched', String(user.watched)).then(() => {
            this.redis.expire(`usersStats:${stream.user_id}:${user.userId}`, 1200);
          });
        }
      }); */
    }));

    console.log(`End ${performance.now() - start}ms`);
  }
}