import { Injectable, Logger } from '@nestjs/common';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';

import { typeorm } from '../../libs/typeorm.js';

@Injectable()
export class TimersService {
  private readonly logger = new Logger(TimersService.name);

  async handleStreamChange(channelId: string) {
    this.logger.log(`Updating timers for streamer: ${channelId}`);
    await typeorm
      .getRepository(ChannelTimer)
      .update({ channelId }, { lastTriggerMessageNumber: 0, last: 0 });
  }
}
