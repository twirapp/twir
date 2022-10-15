import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { CustomVarsRepository, GreetingsRepository, KeywordsRepository } from '@tsuwari/redis';
import { ClientProxy } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';

import { redisSource } from '../../libs/redis.js';
import { typeorm } from '../../libs/typeorm.js';

@Injectable()
export class CacherService implements OnModuleInit {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;
  private readonly logger = new Logger(CacherService.name);

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
    this.updateGreetingsCacheByChannelId(channelId);
    this.updateKeywordsCacheByChannelId(channelId);
    this.updateVariablesCacheByChannelId(channelId);
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

    const repository = redisSource.getRepository(GreetingsRepository);
    for (const greeting of greetings) {
      repository.write(`${greeting.channelId}:${greeting.userId}`, {
        ...greeting,
        processed: false,
      });
    }
  }

  async updateKeywordsCacheByChannelId(channelId: string) {
    const keywords = await typeorm.getRepository(ChannelKeyword).findBy({
      channelId,
    });

    const repository = redisSource.getRepository(KeywordsRepository);

    for (const keyword of keywords) {
      repository.write(`${channelId}:${keyword.id}`, keyword);
    }
  }

  async updateVariablesCacheByChannelId(channelId: string) {
    const variables = await typeorm.getRepository(ChannelCustomvar).findBy({
      channelId,
    });

    const repository = redisSource.getRepository(CustomVarsRepository);

    for (const variable of variables) {
      repository.write(`${channelId}:${variable.id}`, variable);
    }
  }
}
