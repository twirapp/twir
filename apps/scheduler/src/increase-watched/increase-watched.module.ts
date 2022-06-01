import { Module } from '@nestjs/common';

import { RedisService } from '../redis.service.js';
import { IncreaseWatchedService } from './increase-watched.service.js';

@Module({
  providers: [IncreaseWatchedService, RedisService],
})
export class IncreaseWatchedModule { }
