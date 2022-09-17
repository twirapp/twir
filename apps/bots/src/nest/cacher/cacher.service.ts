import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { greetingsSchema, keywordsSchema, RedisORMService } from '@tsuwari/redis';
import { ClientProxy } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';

import { addTimerToQueue, removeTimerFromQueue } from '../../libs/timers.js';
import { typeorm } from '../../libs/typeorm.js';

@Injectable()
export class CacherService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  private readonly logger = new Logger(CacherService.name);

  constructor(private readonly redis: RedisORMService) {}

  async onModuleInit() {
    const channels = await typeorm.getRepository(Channel).find({
      select: {
        id: true,
      },
    });

    for (const channel of channels) {
      this.updateChannel(channel.id);
    }
  }

  async updateChannel(channelId: string) {
    this.logger.log(`Updating systems cache for ${channelId}`);
    this.updateCommandsCacheByChannelId(channelId);
    this.updateTimersCacheByChannelId(channelId);
    this.updateGreetingsCacheByChannelId(channelId);
    this.updateKeywordsCacheByChannelId(channelId);
    this.updateVariablesCacheByChannelId(channelId);
  }

  async updateTimersCacheByChannelId(channelId: string) {
    const timers = await typeorm.getRepository(ChannelTimer).findBy({
      channelId,
    });

    for (const timer of timers) {
      if (timer.enabled) {
        await addTimerToQueue(timer);
      } else {
        removeTimerFromQueue(timer);
      }
    }
  }

  async updateCommandsCacheByChannelId(channelId: string) {
    const commands = await typeorm.getRepository(ChannelCommand).find({
      where: { channelId },
      relations: { responses: true },
    });

    for (const command of commands) {
      await this.nats.send('setCommandCache', command).toPromise();
    }
  }

  async updateGreetingsCacheByChannelId(channelId: string) {
    const greetings = await typeorm.getRepository(ChannelGreeting).findBy({
      channelId,
    });

    const repository = this.redis.fetchRepository(greetingsSchema);
    for (const greeting of greetings) {
      repository.createAndSave(
        {
          ...greeting,
          processed: false,
        },
        `${greeting.channelId}:${greeting.userId}`,
      );
    }
  }

  async updateKeywordsCacheByChannelId(channelId: string) {
    const keywords = await typeorm.getRepository(ChannelKeyword).findBy({
      channelId,
    });

    const repository = this.redis.fetchRepository(keywordsSchema);

    for (const keyword of keywords) {
      repository.createAndSave(keyword, `${channelId}:${keyword.id}`);
    }
  }

  async updateVariablesCacheByChannelId(channelId: string) {
    const variables = await typeorm.getRepository(ChannelCustomvar).findBy({
      channelId,
    });

    for (const variable of variables) {
      this.redis.set(`variables:${channelId}:${variable.id}`, JSON.stringify(variable));
    }
  }
}
