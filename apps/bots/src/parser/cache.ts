import { ChannelIntegration, CustomVar, IntegrationService, UserStats } from '@tsuwari/prisma';
import { CachedStream } from '@tsuwari/shared';
import { getRawData } from '@twurple/common';

import { staticApi } from '../bots.js';
import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { FaceitIntegration } from '../integrations/faceit.js';
import { prisma } from '../libs/prisma.js';
import { redis } from '../libs/redis.js';


export type StatsOfUser = { [key in keyof UserStats]: string } | null;

export class ParserCache {
  #stream: CachedStream | null;
  #userStats: StatsOfUser | null;
  #faceit: Awaited<ReturnType<typeof FaceitIntegration.fetchStats>> | null;
  #enabledIntegrations: (ChannelIntegration & {
    integration: {
      service: IntegrationService;
    };
  })[];
  #customVars: CustomVar[] | null;

  constructor(private readonly broadcasterId: string, private readonly senderId: string) { }

  async getEnabledIntegrations() {
    if (this.#enabledIntegrations) return this.#enabledIntegrations;
    const integrations = await prisma.channelIntegration.findMany({
      where: {
        enabled: true,
      },
      include: {
        integration: {
          select: { service: true },
        },
      },
    });

    this.#enabledIntegrations = integrations;

    return this.#enabledIntegrations;
  }

  async getCustomVars() {
    if (this.#customVars) return this.#customVars;

    const vars = await prisma.customVar.findMany({
      where: {
        channelId: this.broadcasterId,
      },
    });

    this.#customVars = vars;

    return this.#customVars;
  }

  async getStream() {
    if (this.#stream) return this.#stream;

    const streamKey = `streams:${this.broadcasterId}`;
    const cachedStream = await redis.get(streamKey);
    if (cachedStream) {
      this.#stream = JSON.parse(cachedStream) as CachedStream;
    } else {
      const stream = await staticApi.streams.getStreamByUserId(this.broadcasterId);

      if (stream) {
        this.#stream = getRawData(stream);
        redis.set(streamKey, JSON.stringify(this.#stream));
      } else {
        this.#stream = null;
      }
    }

    return this.#stream;
  }

  async getUserStats(): Promise<StatsOfUser | null> {
    if (this.#userStats) {
      return this.#userStats;
    }

    const userKey = `usersStats:${this.broadcasterId}:${this.senderId}`;
    const redisStats = await redis.hgetall(userKey);
    if (!Object.keys(redisStats).length) {
      this.#userStats = await prisma.userStats.findFirst({
        where: {
          userId: this.senderId,
          channelId: this.broadcasterId,
        },
      }) as unknown as StatsOfUser;

      if (this.#userStats) {
        redis.hmset(userKey, this.#userStats).then(() => {
          redis.expire(userKey, USERS_STATUS_CACHE_TTL);
        });
      }
    } else {
      this.#userStats = redisStats as unknown as StatsOfUser;
    }

    return this.#userStats;
  }

  async getFaceitData(nickname: string, game?: string) {
    if (this.#faceit) return this.#faceit;

    const redisKey = `faceit:data:${nickname}`;
    const cachedData = await redis.get(`faceit:data:${nickname}`);

    if (cachedData) {
      this.#faceit = JSON.parse(cachedData);
      return this.#faceit;
    }

    const data = await FaceitIntegration.fetchStats(nickname, game);
    // 5 minutes cache
    redis.set(redisKey, JSON.stringify(data), 'EX', 5 * 60);

    this.#faceit = data as unknown as Awaited<ReturnType<typeof FaceitIntegration.fetchStats>> | null;
    return this.#faceit;
  }
}
