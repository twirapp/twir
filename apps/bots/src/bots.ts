import { config } from '@tsuwari/config';
import { Greeting } from '@tsuwari/prisma';
import { MyRefreshingProvider } from '@tsuwari/shared';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import chunk from 'lodash.chunk';
import pc from 'picocolors';

import { Bot } from './libs/bot.js';
import { prisma } from './libs/prisma.js';
import { redis } from './libs/redis.js';

const staticProvider = new ClientCredentialsAuthProvider(
  config.TWITCH_CLIENTID,
  config.TWITCH_CLIENTSECRET,
);
export const staticApi = new ApiClient({ authProvider: staticProvider });

class BotsClass {
  cache: Map<string, Bot> = new Map();

  async init() {
    const bots = await prisma.bot.findMany({
      include: {
        token: true,
      },
    });

    for (const bot of bots.filter((b) => b.token)) {
      if (!bot.token) continue;
      const botInfo = await staticApi.users.getUserById(bot.id)!;

      if (!botInfo) {
        console.error(`Cannot fetch twitch info for ${bot.id}`);
        continue;
      }

      const authProvider = new MyRefreshingProvider({
        clientId: config.TWITCH_CLIENTID,
        clientSecret: config.TWITCH_CLIENTSECRET,
        prisma,
        token: bot.token,
      });

      const getChannels = async (): Promise<string[]> => {
        const channelsForBot = await prisma.channel.findMany({
          where: {
            isEnabled: true,
            isBanned: false,
            isTwitchBanned: false,
            botId: bot.id,
          },
        });
        const ids = channelsForBot.map((c) => c.id);
        const names = await Promise.all(
          chunk(ids, 100).map(async (chunk) => {
            const channels = await staticApi.users.getUsersByIds(chunk);
            return channels.map((c) => c.name);
          }),
        );

        return names.flat();
      };
      const instance = new Bot(authProvider, getChannels, bot.id);

      console.log(
        `${pc.bgCyan(pc.black('!'))} ${pc.magenta(botInfo.name)} ${pc.bgYellow(
          'Connecting to twitch...',
        )}`,
      );
      await instance.connect();

      this.cache.set(bot.id, instance);
    }
  }
}

export const Bots = new BotsClass();
