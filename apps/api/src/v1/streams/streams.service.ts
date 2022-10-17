import { Injectable } from '@nestjs/common';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';

import { typeorm } from '../../index.js';

@Injectable()
export class StreamsService {
  async #resetGreetings(channelId: string) {
    await typeorm.getRepository(ChannelGreeting).update({ channelId }, { processed: false });
  }

  async handleStreamStateChange(channelId: string) {
    this.#resetGreetings(channelId);
  }

  async getStream(channelId: string) {
    const stream = await typeorm.getRepository(ChannelStream).findOneBy({
      userId: channelId,
    });
    if (!stream) return null;
    return stream;
  }
}
