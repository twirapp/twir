import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { RedisORMService, greetingsSchema, keywordsSchema, customVarSchema } from '@tsuwari/redis';
import { ClientProxy } from '@tsuwari/shared';

import { prisma } from '../../libs/prisma.js';
import { removeTimerFromQueue, addTimerToQueue } from '../../libs/timers.js';

@Injectable()
export class CacherService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  private readonly logger = new Logger(CacherService.name);

  constructor(private readonly redis: RedisORMService) { }

  async onModuleInit() {
    const channels = await prisma.channel.findMany({
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
    const timers = await prisma.timer.findMany({
      where: {
        channelId,
      },
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
    const commands = await prisma.command.findMany({
      where: {
        channelId,
      },
      include: {
        responses: true,
      },
    });

    for (const command of commands) {
      await this.nats.send('setCommandCache', command).toPromise();
    }
  }

  async updateGreetingsCacheByChannelId(channelId: string) {
    const greetings = await prisma.greeting.findMany({
      where: {
        channelId,
      },
    });

    const repository = this.redis.fetchRepository(greetingsSchema);
    for (const greeting of greetings) {
      repository.createAndSave({
        ...greeting,
        processed: false,
      }, `${greeting.channelId}:${greeting.userId}`);
    }
  }

  async updateKeywordsCacheByChannelId(channelId: string) {
    const keywords = await prisma.keyword.findMany({
      where: {
        channelId,
      },
    });

    const repository = this.redis.fetchRepository(keywordsSchema);

    for (const keyword of keywords) {
      repository.createAndSave(keyword, `${channelId}:${keyword.id}`);
    }
  }

  async updateVariablesCacheByChannelId(channelId: string) {
    const variables = await prisma.customVar.findMany({
      where: {
        channelId,
      },
    });

    const repository = this.redis.fetchRepository(customVarSchema);

    for (const variable of variables) {
      repository.createAndSave(variable, `${channelId}:${variable.id}`);
    }
  }
}
