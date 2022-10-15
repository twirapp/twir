import 'reflect-metadata';

import { Global, Injectable, Module, OnModuleInit } from '@nestjs/common';
import { config } from '@tsuwari/config';
import IORedis from 'ioredis';
import { Client, Schema } from 'redis-om';

import { BaseRepository } from './base.js';
import * as repositories from './repositories/index.js';
import * as schemas from './schemas/index.js';

export { Client, Repository } from 'redis-om';
export * from './repositories/index.js';
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

export class RedisSource {
  repositories: Map<string, BaseRepository<any>> = new Map();
  private redis: IORedis;

  constructor(redis?: IORedis) {
    this.redis = redis ?? new IORedis(config.REDIS_URL);
    for (const Repo of Object.values(repositories)) {
      this.repositories.set(Repo.constructor.name, new Repo(this.redis));
    }
  }

  getRepository<T extends abstract new (...args: any) => any>(t: T) {
    const repo = this.repositories.get((t as any).constuctor.name) as InstanceType<T> | undefined;
    if (!repo) throw new Error(`Seems like ${t.constructor.name} not registered in repositories.`);
    return repo;
  }
}
