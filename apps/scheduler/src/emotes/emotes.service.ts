import { Injectable } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { config } from '@tsuwari/config';
import { Channel } from '@tsuwari/typeorm/entities/Channel';

import { typeorm } from '../index.js';
import { emotesCacherGrpcClient } from '../libs/emotes.grpc.js';

@Injectable()
export class EmotesService {

  @Interval('cacheEmotes', config.isDev ? 5000 : 5 * 60 * 1000)
  async cacheEmotes() {
    const channels = await typeorm.getRepository(Channel).find({
      select: {
        id: true,
      },
    });

    for (const channel of channels) {
      emotesCacherGrpcClient.cacheEmotes({ channelId: channel.id });
    }
  }
}