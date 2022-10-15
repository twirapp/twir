import type { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

type Greeting = Omit<ChannelGreeting, 'channel'> & {
  processed: boolean;
};

export class GreetingsRepository extends BaseRepository<Greeting> {
  constructor(redis: Redis) {
    super('greetings', redis);
  }
}
