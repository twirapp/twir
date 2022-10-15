import type { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

export class GreetingsRepository extends BaseRepository<Omit<ChannelGreeting, 'channel'>> {
  constructor(redis: Redis) {
    super('greetings', redis);
  }
}
