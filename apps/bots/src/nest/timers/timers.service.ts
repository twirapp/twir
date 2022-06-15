import { Injectable, Logger } from '@nestjs/common';

import { prisma } from '../../libs/prisma.js';

@Injectable()
export class TimersService {
  private readonly logger = new Logger(TimersService.name);

  async handleStreamChange(channelId: string) {
    this.logger.log(`Updating timers for streamer: ${channelId}`);
    await prisma.timer.updateMany({
      where: {
        channelId,
      },
      data: {
        lastTriggerMessageNumber: 0,
        last: 0,
      },
    });
  }
}
