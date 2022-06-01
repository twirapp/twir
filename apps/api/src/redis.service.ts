import { Global, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';
import Redis from 'ioredis';

@Global()
@Injectable()
export class RedisService extends Redis {
  constructor() {
    super(config.REDIS_URL);
  }
}
