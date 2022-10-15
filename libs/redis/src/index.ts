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

export class RedisSource {
  #repositories: Map<string, BaseRepository<any>> = new Map();
  #redis: IORedis;

  constructor(redis?: IORedis) {
    this.#redis = redis ?? new IORedis(config.REDIS_URL);
    const reps = [
      repositories.CustomVarsRepository,
      repositories.GreetingsRepository,
      repositories.KeywordsRepository,
      repositories.ModerationSettingsRepository,
      repositories.StreamRepository,
      repositories.UsersStatsRepository,
    ];
    for (const Repo of reps) {
      this.#repositories.set(Repo.name, new Repo(this.#redis));
    }
  }

  getRepository<T extends abstract new (...args: any) => any>(t: T) {
    const repo = this.#repositories.get((t as any).name) as InstanceType<T> | undefined;
    if (!repo) throw new Error(`Seems like ${t.name} not registered in repositories.`);
    return repo;
  }
}

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
@Injectable()
export class RedisDataSourceService extends RedisSource {}

@Global()
@Module({
  controllers: [],
  providers: [RedisORMService, RedisDataSourceService],
  exports: [RedisORMService, RedisDataSourceService],
})
export class RedisORMModule {}
