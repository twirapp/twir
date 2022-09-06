import { Injectable, OnModuleInit } from '@nestjs/common';
import { Greetings, greetingsSchema, RedisORMService, Repository } from '@tsuwari/redis';
import { RedisService } from '@tsuwari/shared';

@Injectable()
export class StreamsService implements OnModuleInit {
  #greetingsRepository: Repository<Greetings>;

  constructor(private readonly redis: RedisService, private readonly redisOrm: RedisORMService) {}

  onModuleInit() {
    this.#greetingsRepository = this.redisOrm.fetchRepository(greetingsSchema);
  }

  async #resetGreetings(channelId: string) {
    const greetings = await this.#greetingsRepository
      .search()
      .where('channelId')
      .equal(channelId)
      .returnAll();

    for (const greeting of greetings.map((g) => g.toRedisJson())) {
      this.#greetingsRepository.createAndSave(
        {
          ...greeting,
          processed: false,
        },
        `${greeting.channelId}:${greeting.userId}`,
      );
    }
  }

  async handleStreamStateChange(channelId: string) {
    this.#resetGreetings(channelId);
  }

  async getStream(channelId: string) {
    const stream = await this.redis.get(`streams:${channelId}`);
    if (!stream) return null;
    return JSON.parse(stream);
  }
}
