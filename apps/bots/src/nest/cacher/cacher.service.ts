import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy } from '@tsuwari/shared';

import { prisma } from '../../libs/prisma.js';
import { redis } from '../../libs/redis.js';
import { removeTimerFromQueue, addTimerToQueue } from '../../libs/timers.js';

@Injectable()
export class CacherService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  private readonly logger = new Logger(CacherService.name);

  async onModuleInit() {
    const channels = await prisma.channel.findMany({
      select: {
        id: true,
      },
    });

    for (const channel of channels) {
      this.logger.log(`Updating systems cache for ${channel.id}`);
      this.updateCommandsCacheByChannelId(channel.id);
      this.updateTimersCacheByChannelId(channel.id);
      this.updateGreetingsCacheByChannelId(channel.id);
      this.updateKeywordsCacheByChannelId(channel.id);
      this.updateVariablesCacheByChannelId(channel.id);
    }
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
        await removeTimerFromQueue(timer);
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

    for (const greeting of greetings) {
      await redis.hset(`greetings:${greeting.channelId}:${greeting.userId}`, greeting);
    }
  }

  async updateKeywordsCacheByChannelId(channelId: string) {
    const keywords = await prisma.keyword.findMany({
      where: {
        channelId,
      },
    });

    for (const keyword of keywords) {
      await redis.hmset(`keywords:${channelId}:${keyword.id}`, keyword);
    }
  }

  async updateVariablesCacheByChannelId(channelId: string) {
    const variables = await prisma.customVar.findMany({
      where: {
        channelId,
      },
    });

    for (const variable of variables) {
      await redis.set(`variables:${channelId}:${variable.name}`, JSON.stringify(variable));
    }
  }
}
