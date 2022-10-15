import type { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import Redis from 'ioredis';

import { BaseRepository } from '../base.js';

export class ModerationSettingsRepository extends BaseRepository<
  Omit<ChannelModerationSetting, 'channel'>
> {
  constructor(redis: Redis) {
    super('v', redis);
  }
}
