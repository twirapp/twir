import { Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';

import { RedisService } from '../../redis.service.js';

@Injectable()
export class StreamsService {
  constructor(private readonly redis: RedisService, private readonly prisma: PrismaService) { }

  async #resetGreetings(channelId: string) {
    const keys = await this.redis.keys(`greetings:${channelId}:*`);

    for (const key of keys) {
      const greeting = await this.redis.hgetall(key);
      if (!Object.keys(greeting).length) continue;

      await this.redis.hset(`greetings:${greeting.channelId}:${greeting.userId}`, {
        ...greeting,
        processed: false,
      });
    }
  }

  async handleStreamStateChange(channelId: string) {
    this.#resetGreetings(channelId);
  }

  async getStream(channelId: string) {
    const stream = await this.redis.get(`streams:${channelId}`);

    return stream;
  }
}
