import { Injectable, OnModuleInit } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import { Greetings, greetingsSchema, RedisORMService, Repository, Stream, streamSchema } from '@tsuwari/redis';
import { RedisService } from '@tsuwari/shared';


@Injectable()
export class StreamsService implements OnModuleInit {
  #greetingsRepository: Repository<Greetings>;
  #streamRepository: Repository<Stream>;

  constructor(
    private readonly redis: RedisService,
    private readonly prisma: PrismaService,
    private readonly redisOrm: RedisORMService,
  ) { }

  onModuleInit() {
    this.#greetingsRepository = this.redisOrm.fetchRepository(greetingsSchema);
    this.#streamRepository = this.redisOrm.fetchRepository(streamSchema);
  }

  async #resetGreetings(channelId: string) {
    const greetings = await this.#greetingsRepository.search()
      .where('channelId').equal(channelId)
      .returnAll();

    for (const greeting of greetings.map(g => g.toRedisJson())) {
      this.#greetingsRepository.createAndSave({
        ...greeting,
        processed: false,
      }, `${greeting.channelId}:${greeting.userId}`);
    }
  }

  async handleStreamStateChange(channelId: string) {
    this.#resetGreetings(channelId);
  }

  async getStream(channelId: string) {
    const s = await this.#streamRepository.fetch(channelId);
    return s.toRedisJson();
  }
}
