import { Injectable } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';
import _ from 'lodash';

import { typeorm } from '../index.js';
import { watchedGrpcClient } from '../libs/watched.grpc.js';

@Injectable()
export class WatchedService {
  @Interval(config.isDev ? 5 * 1000 : 5 * 60 * 1000)
  async watched() {
    const streams = await typeorm.getRepository(ChannelStream).find({
      where: {
        channel: { isEnabled: true },
      },
      select: {
        userId: true,
        channel: {
          botId: true,
          isEnabled: true,
        },
      },
      relations: {
        channel: true,
      },
    });

    const groups = _.groupBy(streams, (s) => s.channel!.botId!);

    for (const [botId, channels] of Object.entries(groups)) {
      const chunks = _.chunk(channels, 100);

      for (const ch of chunks) {
        const mapped = ch.map((c) => c.userId);
        watchedGrpcClient.incrementByChannelId({
          botId: botId,
          channelsId: mapped,
        });
      }
    }
  }
}
