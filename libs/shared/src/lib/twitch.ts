import { PrismaClient, Token } from '@tsuwari/prisma';
import { RefreshingAuthProvider } from '@twurple/auth';

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