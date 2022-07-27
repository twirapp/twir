import { ChannelIntegration, IntegrationService, PrismaService } from '@tsuwari/prisma';
import { RedisORMService, CustomVar, UsersStats, usersStatsSchema, customVarSchema, streamSchema, Stream } from '@tsuwari/redis';
import { RedisService, TwitchApiService } from '@tsuwari/shared';
import { getRawData } from '@twurple/common';

import { USERS_STATUS_CACHE_TTL } from '../constants.js';
import { FaceitIntegration } from '../integrations/faceit.js';

export class ParserCache {
  #stream: ReturnType<typeof Stream.prototype.toRedisJson> | null;
  #userStats: ReturnType<typeof UsersStats.prototype.toRedisJson>;
  #faceit: Awaited<ReturnType<typeof FaceitIntegration.prototype.fetchStats>> | null;
  #enabledIntegrations: (ChannelIntegration & {
    integration: {
      service: IntegrationService;
    };
  })[];
  #customVars: ReturnType<typeof CustomVar.prototype.toRedisJson>[] | null;

  constructor(
    private readonly staticApi: TwitchApiService,
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
    private readonly faceitIntegration: FaceitIntegration,
    private readonly broadcasterId: string,
    private readonly redisOm: RedisORMService,
    private readonly senderId?: string,
  ) { }

  async getEnabledIntegrations(channelId: string) {
    if (this.#enabledIntegrations) return this.#enabledIntegrations;
    const integrations = await this.prisma.channelIntegration.findMany({
      where: {
        enabled: true,
        channelId,
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

    const repository = this.redisOm.fetchRepository(customVarSchema);
    const redisVars = await repository.search()
      .where('channelId').equals(this.broadcasterId)
      .returnAll()
      .then(data => {
        return data.map(v => v.toRedisJson()).filter(v => Object.keys(v).length);
      });

    if (redisVars.length) {
      this.#customVars = redisVars;
      return this.#customVars;
    }

    const vars = await this.prisma.customVar.findMany({
      where: {
        channelId: this.broadcasterId,
      },
    });

    for (const variable of vars) {
      repository.createAndSave(variable, `${this.broadcasterId}:${variable.name}`);
    }

    this.#customVars = vars;

    return this.#customVars;
  }

  async getStream() {
    if (this.#stream) return this.#stream;

    const repository = this.redisOm.fetchRepository(streamSchema);

    const cachedStream = await repository.fetch(this.broadcasterId).then(s => s.toRedisJson());
    if (Object.keys(cachedStream).length) {
      this.#stream = cachedStream;
    } else {
      const stream = await this.staticApi.streams.getStreamByUserId(this.broadcasterId);

      if (stream) {
        this.#stream = getRawData(stream);
        repository.createAndSave(this.#stream, this.broadcasterId);
      } else {
        this.#stream = null;
      }
    }

    return this.#stream;
  }

  async getUserStats() {
    if (this.#userStats) {
      return this.#userStats;
    }
    const repository = this.redisOm.fetchRepository(usersStatsSchema);

    const userKey = `${this.broadcasterId}:${this.senderId}`;
    const redisStats = await repository.fetch(userKey).then(r => r.toRedisJson());

    if (!Object.keys(redisStats).length) {
      const stats = await this.prisma.userStats.findFirst({
        where: {
          userId: this.senderId,
          channelId: this.broadcasterId,
        },
      });

      if (stats) {
        repository.createAndSave({
          ...stats,
          watched: stats.watched.toString(),
        }, userKey).then(() => {
          repository.expire(userKey, USERS_STATUS_CACHE_TTL);
        });
        this.#userStats = {
          ...stats,
          watched: stats.watched.toString(),
        };
      }
    } else {
      this.#userStats = redisStats;
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
      console.log(e);
      return null;
    }
  }
}
