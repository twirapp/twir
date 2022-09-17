import { CachedGetter } from '@d-fischer/cache-decorators';
import { Global, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { Repository } from '@tsuwari/typeorm';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { ApiClient, HelixUserApi, HelixUserData, UserIdResolvable } from '@twurple/api';
import { ClientCredentialsAuthProvider, RefreshingAuthProvider } from '@twurple/auth';
import { getRawData, UserNameResolvable } from '@twurple/common';
import Redis from 'ioredis';

import { RedisService } from './redis.js';
import { WEEK } from './time.js';

export class MyRefreshingProvider extends RefreshingAuthProvider {
  constructor(opts: {
    clientId: string;
    clientSecret: string;
    repository: Repository<Token>;
    token: Token;
  }) {
    super(
      {
        clientId: opts.clientId,
        clientSecret: opts.clientSecret,
        onRefresh: async (refreshedToken) => {
          const { accessToken, refreshToken, obtainmentTimestamp, expiresIn } = refreshedToken;
          if (!refreshToken || !obtainmentTimestamp || !expiresIn) return;
          await opts.repository.update(
            { id: opts.token.id },
            {
              accessToken,
              refreshToken,
              obtainmentTimestamp: new Date(obtainmentTimestamp),
              expiresIn,
            },
          );
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

  async getUserByNameWithCache(userName: UserNameResolvable): Promise<HelixUserData> {
    const redisKey = `twitchUsersCache:${userName}`;
    let data: HelixUserData | null = null;

    const cachedData = await this.redis?.get(redisKey);
    if (cachedData) {
      data = JSON.parse(cachedData) as HelixUserData;
    } else {
      const user = await super.getUserByName(userName);
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
