import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { ClientProxy } from '@tsuwari/shared';
import { Channel } from '@tsuwari/typeorm/entities/Channel';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';

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
    this.updateGreetingsCacheByChannelId(channelId);
  }

  async updateGreetingsCacheByChannelId(channelId: string) {
    await typeorm.getRepository(ChannelGreeting).update({ channelId }, { processed: false });
  }
}
