import { Injectable } from '@nestjs/common';
import {
  GreetingsRepository,
  RedisDataSourceService,
  StreamRepository,
  StreamType,
} from '@tsuwari/redis';
import { RedisService } from '@tsuwari/shared';

@Injectable()
export class StreamsService {
  constructor(
    private readonly redis: RedisService,
    private readonly redisSource: RedisDataSourceService,
  ) {}

  async #resetGreetings(channelId: string) {
    const repository = this.redisSource.getRepository(GreetingsRepository);

    const greetingsKeys = await this.redis.keys(`greetings:${channelId}:*`);
    const greetings = await repository.readMany(greetingsKeys, true);

    for (const greeting of greetings) {
      repository.write(`${greeting.channelId}:${greeting.userId}`, {
        ...greeting,
        processed: false,
      });
    }
  }

  async handleStreamStateChange(channelId: string) {
    this.#resetGreetings(channelId);
  }

  async getStream(channelId: string): Promise<StreamType | null> {
    const stream = await this.redisSource.getRepository(StreamRepository).read(channelId);
    return stream;
  }
}
