import { Controller } from '@nestjs/common';
import { EventPattern } from '@nestjs/microservices';
import { ClientProxyEvents } from '@tsuwari/shared';

import { addTimerToQueue, removeTimerFromQueue } from '../../libs/timers.js';
import { TimersService } from './timers.service.js';

@Controller()
export class TimersController {
  constructor(private readonly service: TimersService) { }

  @EventPattern('bots.addTimerToQueue')
  addTimerToQueue(timerId: string) {
    addTimerToQueue(timerId);
  }

  @EventPattern('bots.removeTimerFromQueue')
  removeTimerFromQueue(timerId: string) {
    removeTimerFromQueue(timerId);
  }

  @EventPattern('streams.online')
  streamOnline(data: ClientProxyEvents['streams.online']['input']) {
    this.service.handleStreamChange(data.channelId);
  }

  @EventPattern('streams.offline')
  streamOffline(data: ClientProxyEvents['streams.offline']['input']) {
    this.service.handleStreamChange(data.channelId);
  }
}
