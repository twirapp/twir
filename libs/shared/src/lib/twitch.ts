import { CachedGetter } from '@d-fischer/cache-decorators';
import { Global, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaClient, Token } from '@tsuwari/prisma';
import { ApiClient, HelixUserApi, HelixUserData, UserIdResolvable } from '@twurple/api';
import { ClientCredentialsAuthProvider, RefreshingAuthProvider } from '@twurple/auth';
import { getRawData } from '@twurple/common';
import Redis from 'ioredis';

import { RedisService } from './redis.js';
import { WEEK } from './time.js';

export class MyRefreshingProvider extends RefreshingAuthProvider {
  constructor(opts: {
    clientId: string;
    clientSecret: string;
    prisma: PrismaClient;
    token: Token;
  }) {
    super(
      {
        clientId: opts.clientId,
        clientSecret: opts.clientSecret,
        onRefresh: async (refreshedToken) => {
          const { accessToken, refreshToken, obtainmentTimestamp, expiresIn } = refreshedToken;
          if (!refreshToken || !obtainmentTimestamp || !expiresIn) return;
          await opts.prisma.token.update({
            where: {
              id: opts.token.id,
            },
            data: {
              accessToken,
              refreshToken,
              obtainmentTimestamp: new Date(obtainmentTimestamp),
              expiresIn,
            },
          });
        },
      },
      {
        refreshToken: opts.token.refreshToken,
        expiresIn: opts.token.expiresIn,
        obtainmentTimestamp: opts.token.obtainmentTimestamp.getTime(),
      },
    );
  }
}

class MyUserApi extends HelixUserApi {
  constructor(client: TwitchApiService, readonly redis?: RedisService) {
    super(client);
  }

  async getUserByIdWithCache(userId: UserIdResolvable): Promise<HelixUserData> {
    const redisKey = `twitchUsersCache:${userId}`;
    let data: HelixUserData | null = null;

    const cachedData = await this.redis?.get(redisKey);
    if (cachedData) {
      data = JSON.parse(cachedData) as HelixUserData;
    } else {
      const user = await super.getUserById(userId);
      if (user) {
        data = getRawData(user);
        this.redis?.set(redisKey, JSON.stringify(data), 'EX', (WEEK * 2) / 1000);
      }
    }

    return data;
  }
}

@Global()
@Injectable()
export class TwitchApiService extends ApiClient {
  constructor(readonly redis?: RedisService) {
    const staticProvider = new ClientCredentialsAuthProvider(
      config.TWITCH_CLIENTID,
      config.TWITCH_CLIENTSECRET,
    );
    super({
      authProvider: staticProvider,
    });
  }

  @CachedGetter()
  get users() {
    return new MyUserApi(this, this.redis);
  }
}
