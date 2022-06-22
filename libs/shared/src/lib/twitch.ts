import { Global, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { PrismaClient, Token } from '@tsuwari/prisma';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider, RefreshingAuthProvider } from '@twurple/auth';

export class MyRefreshingProvider extends RefreshingAuthProvider {
  constructor(opts: {
    clientId: string,
    clientSecret: string,
    prisma: PrismaClient,
    token: Token
  }) {
    super({
      clientId: opts.clientId,
      clientSecret: opts.clientSecret,
      onRefresh: async (refreshedToken) => {
        const { accessToken, refreshToken, obtainmentTimestamp, expiresIn } = refreshedToken;
        if (!refreshToken || !obtainmentTimestamp || !expiresIn) return;
        await opts.prisma.token.update({
          where: {
            id: opts.token.id,
          },
          data: { accessToken, refreshToken, obtainmentTimestamp: new Date(obtainmentTimestamp), expiresIn },
        });
      },
    }, { refreshToken: opts.token.refreshToken, expiresIn: opts.token.expiresIn, obtainmentTimestamp: opts.token.obtainmentTimestamp.getTime() });
  }
}

@Global()
@Injectable()
export class TwitchApiService extends ApiClient {
  constructor() {
    const staticProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
    super({
      authProvider: staticProvider,
    });
  }
}