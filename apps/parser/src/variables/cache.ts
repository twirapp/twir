import { ChannelIntegration, CustomVar, IntegrationService, PrismaService, UserStats } from '@tsuwari/prisma';
import { CachedStream, RedisService, TwitchApiService } from '@tsuwari/shared';
import { getRawData } from '@twurple/common';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { FaceitIntegration } from '../integrations/faceit.js';

export type StatsOfUser = { [key in keyof UserStats]: string } | null;

export class ParserCache {
  #stream: CachedStream | null;
  #userStats: StatsOfUser | null;
  #faceit: Awaited<ReturnType<typeof FaceitIntegration.prototype.fetchStats>> | null;
  #enabledIntegrations: (ChannelIntegration & {
    integration: {
      service: IntegrationService;
    };
  })[];
  #customVars: CustomVar[] | null;

  constructor(
    private readonly staticApi: TwitchApiService,
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly faceitIntegration: FaceitIntegration,
    private readonly broadcasterId: string,
    private readonly senderId?: string,
  ) { }

  async getEnabledIntegrations() {
    if (this.#enabledIntegrations) return this.#enabledIntegrations;
    const integrations = await this.prisma.channelIntegration.findMany({
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

    const keys = await this.redis.keys(`variables:${this.broadcasterId}:*`);

    if (keys.length) {
      const variables = await Promise.all(keys.map(k => this.redis.get(k)));
      this.#customVars = variables.map(v => JSON.parse(v!) as CustomVar);
      return this.#customVars;
    }

    const vars = await this.prisma.customVar.findMany({
      where: {
        channelId: this.broadcasterId,
      },
    });

    for (const variable of vars) {
      this.redis.set(`variables:${this.broadcasterId}:${variable.name}`, JSON.stringify(variable));
    }

    this.#customVars = vars;

    return this.#customVars;
  }

  async getStream() {
    if (this.#stream) return this.#stream;

    const streamKey = `streams:${this.broadcasterId}`;
    const cachedStream = await this.redis.get(streamKey);
    if (cachedStream) {
      this.#stream = JSON.parse(cachedStream) as CachedStream;
    } else {
      const stream = await this.staticApi.streams.getStreamByUserId(this.broadcasterId);

      if (stream) {
        this.#stream = getRawData(stream);
        this.redis.set(streamKey, JSON.stringify(this.#stream));
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
    const redisStats = await this.redis.hgetall(userKey);
    if (!Object.keys(redisStats).length) {
      this.#userStats = await this.prisma.userStats.findFirst({
        where: {
          userId: this.senderId,
          channelId: this.broadcasterId,
        },
      }) as unknown as StatsOfUser;

      if (this.#userStats) {
        this.redis.hmset(userKey, this.#userStats).then(() => {
          this.redis.expire(userKey, USERS_STATUS_CACHE_TTL);
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
    const cachedData = await this.redis.get(`faceit:data:${nickname}`);

    if (cachedData) {
      this.#faceit = JSON.parse(cachedData);
      return this.#faceit;
    }

    try {
      const data = await this.faceitIntegration.fetchStats(nickname, game);
      // 5 minutes cache
      this.redis.set(redisKey, JSON.stringify(data), 'EX', 5 * 60);

      this.#faceit = data as unknown as Awaited<ReturnType<typeof FaceitIntegration.prototype.fetchStats>> | null;
      return this.#faceit;
    } catch (e) {
      return null;
    }
  }
}
