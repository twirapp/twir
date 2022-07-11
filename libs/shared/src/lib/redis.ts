import { Global, Injectable, Module } from '@nestjs/common';
import { config } from '@tsuwari/config';
import Redis from 'ioredis';

@Global()
@Injectable()
export class RedisService extends Redis {
  constructor() {
    super(config.REDIS_URL);
  }
}

@Global()
@Module({
  controllers: [],
  providers: [RedisService],
  exports: [RedisService],
})
export class RedisModule { }