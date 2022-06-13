import { MyRefreshingProvider } from '@tsuwari/shared';
import { ApiClient } from '@twurple/api';

import { config } from '../../config.js';
import { prisma } from '../../libs/prisma.js';
import { DefaultCommand } from '../types.js';

export const channelInfo: DefaultCommand[] = [
  {
    name: 'game set',
    description: 'Changing game of the channel.',
    permission: 'MODERATOR',
    async handler(state, params?) {
      if (!state.channelId) return;
      if (!params || !params.length) return 'you must specify what game to set.';

      const streamerToken = await prisma.token.findFirst({
        where: {
          user: {
            id: state.channelId,
          },
        },
      });

      if (!streamerToken) return;

      const authProvider = new MyRefreshingProvider({
        clientId: config.TWITCH_CLIENTID,
        clientSecret: config.TWITCH_CLIENTSECRET,
        token: streamerToken,
        prisma,
      });

      const api = new ApiClient({ authProvider });

      const [tokenInfo, game] = await Promise.all([
        api.getTokenInfo(),
        api.games.getGameByName(params),
      ]);

      if (!tokenInfo.scopes.includes('channel:manage:broadcast')) {
        return `Missed scope from streamer. @${state.target.value.substring(1)} should re-login to the dashboard.`;
      }

      if (!game) {
        return `game ${params} not found on twitch.`;
      }

      await api.channels.updateChannelInfo(state.channelId, {
        gameId: game.id,
      });

      return '✅';
    },
  },
  {
    name: 'title set',
    description: 'Changing title of the channel.',
    permission: 'MODERATOR',
    async handler(state, params?) {
      if (!state.channelId) return;
      if (!params || !params.length) return 'you must specify what title to set.';

      const streamerToken = await prisma.token.findFirst({
        where: {
          user: {
            id: state.channelId,
          },
        },
      });

      if (!streamerToken) return;

      const authProvider = new MyRefreshingProvider({
        clientId: config.TWITCH_CLIENTID,
        clientSecret: config.TWITCH_CLIENTSECRET,
        token: streamerToken,
        prisma,
      });

      const api = new ApiClient({ authProvider });

      const tokenInfo = await api.getTokenInfo();
      if (!tokenInfo.scopes.includes('channel:manage:broadcast')) {
        return `Missed scope from streamer. @${state.target.value.substring(1)} should re-login to the dashboard.`;
      }

      await api.channels.updateChannelInfo(state.channelId, {
        title: params,
      });

      return '✅';
    },
  },
];