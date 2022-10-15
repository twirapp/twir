import type { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

export class KeywordsRepository extends BaseRepository<Omit<ChannelKeyword, 'channel'>> {
  constructor(redis: Redis) {
    super('keywords', redis);
  }
}
