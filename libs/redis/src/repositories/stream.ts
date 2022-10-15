import { HelixStreamData } from '@twurple/api';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

export type StreamType = HelixStreamData & { parsedMessages: number };

export class StreamRepository extends BaseRepository<StreamType> {
  constructor(redis: Redis) {
    super('streams', redis);
  }
}
