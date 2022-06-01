import { PrismaClient, Token } from '@tsuwari/prisma';
import { RefreshingAuthProvider } from '@twurple/auth';

export class MyRefreshingProvider extends RefreshingAuthProvider {
  constructor(opts: {
    clientId: string,
    clientSecret: string,
    prisma: PrismaClient,
    token: Token
  }) {
    if (!opts.token.userId && !opts.token.botId) {
      throw new Error('UserId or TokenId should be passed');
    }

    super({
      clientId: opts.clientId,
      clientSecret: opts.clientSecret,
      onRefresh: async (refreshedToken) => {
        const { accessToken, refreshToken, obtainmentTimestamp, expiresIn } = refreshedToken;
        if (!refreshToken || !obtainmentTimestamp || !expiresIn) return;
        await opts.prisma.token.update({
          where: {
            botId: opts.token.botId ?? undefined,
            userId: opts.token.userId ?? undefined,
          },
          data: { accessToken, refreshToken, obtainmentTimestamp: new Date(obtainmentTimestamp), expiresIn },
        });
      },
    }, { refreshToken: opts.token.refreshToken, expiresIn: opts.token.expiresIn, obtainmentTimestamp: opts.token.obtainmentTimestamp.getTime() });
  }
}