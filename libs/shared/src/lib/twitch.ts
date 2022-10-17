import { CachedGetter } from '@d-fischer/cache-decorators';
import { Global, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { Repository } from '@tsuwari/typeorm';
import { Token } from '@tsuwari/typeorm/entities/Token';
import { ApiClient, HelixUserApi, HelixUserData, UserIdResolvable } from '@twurple/api';
import { ClientCredentialsAuthProvider, RefreshingAuthProvider } from '@twurple/auth';
import { getRawData, UserNameResolvable } from '@twurple/common';

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

@Global()
@Injectable()
export class TwitchApiService extends ApiClient {
  constructor() {
    const staticProvider = new ClientCredentialsAuthProvider(
      config.TWITCH_CLIENTID,
      config.TWITCH_CLIENTSECRET,
    );
    super({
      authProvider: staticProvider,
    });
  }
}
