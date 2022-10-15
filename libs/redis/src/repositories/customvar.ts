import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

export class CustomVarsRepository extends BaseRepository<Omit<ChannelCustomvar, 'channel'>> {
  constructor(redis: Redis) {
    super('variables', redis);
  }
}
