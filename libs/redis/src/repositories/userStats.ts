import type { UserStats } from '@tsuwari/typeorm/entities/UserStats';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

type Stats = Omit<UserStats, 'channel' | 'user' | 'watched'> & {
  watched: string;
};

export class UsersStatsRepository extends BaseRepository<Stats> {
  constructor(redis: Redis) {
    super('usersStats', redis);
  }
}
