import type { UserStats } from '@tsuwari/typeorm/entities/UserStats';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

export class UsersStatsRepository extends BaseRepository<Omit<UserStats, 'channel' | 'user'>> {
  constructor(redis: Redis) {
    super('usersStats', redis);
  }
}
