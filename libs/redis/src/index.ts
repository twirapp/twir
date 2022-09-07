import 'reflect-metadata';

import { Global, Injectable, Module, OnModuleInit } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { Client, Schema } from 'redis-om';

import * as schemas from './schemas/index.js';

export { Client, Repository } from 'redis-om';
export * from './schemas/index.js';

@Global()
@Injectable()
export class RedisORMService extends Client implements OnModuleInit {
  async onModuleInit() {
    await this.open(config.REDIS_URL);

    for (const schema of Object.values(schemas).filter((v) => v instanceof Schema)) {
      await this.fetchRepository(schema as any).createIndex();
    }
  }
}

@Global()
@Module({
  controllers: [],
  providers: [RedisORMService],
  exports: [RedisORMService],
})
export class RedisORMModule {}
