import { config } from '@tsuwari/config';
import { Command, Response } from '@tsuwari/prisma';
import { MyRefreshingProvider } from '@tsuwari/shared';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';
import pc from 'picocolors';

import { Bot } from './libs/bot.js';
import { prisma } from './libs/prisma.js';
import { redis } from './libs/redis.js';

const staticProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
export const staticApi = new ApiClient({ authProvider: staticProvider });

class BotsClass {
  cache: Map<string, Bot> = new Map();

  async init() {
    const bots = await prisma.bot.findMany({
      include: {
        token: true,
        channels: {
          where: {
            isEnabled: true,
            isBanned: false,
            isTwitchBanned: false,
          },
        },
      },
    });

    for (const bot of bots.filter((b) => b.token)) {
      if (!bot.token) continue;
      const botInfo = await staticApi.users.getUserById(bot.id)!;

      if (!botInfo) {
        console.error(`Cannot fetch twitch info for ${bot.id}`);
        continue;
      }

      const channels = await staticApi.users.getUsersByIds(bot.channels.map((c) => c.id));

      const authProvider = new MyRefreshingProvider({
        clientId: config.TWITCH_CLIENTID,
        clientSecret: config.TWITCH_CLIENTSECRET,
        prisma,
        token: bot.token,
      });

      const instance = new Bot(
        authProvider,
        channels.map((c) => c.name),
      );

      console.log(`${pc.bgCyan(pc.black('!'))} ${pc.magenta(botInfo.name)} ${pc.bgYellow('Connecting to twitch...')}`);
      await instance.connect();

      this.cache.set(bot.id, instance);

      for (const channel of bot.channels) {
        this.updateGreetingsCacheByChannelid(channel.id);
      }
    }
  }



  async updateGreetingsCacheByChannelid(channelId: string) {
    const greetings = await prisma.greeting.findMany({
      where: { channelId },
    });

    for (const greeting of greetings) {
      await redis.hset(`greetings:${greeting.channelId}:${greeting.userId}`, {
        ...greeting,
        processed: false,
      });
    }

    const keys = await redis.keys(`greetings:${channelId}:*`);

    for (const key of keys) {
      const [, channelId, userId] = key.split(':');
      if (!greetings.some((g) => g.channelId === channelId && g.userId === userId)) {
        await redis.del(key);
      }
    }
  }
}

export const Bots = new BotsClass();