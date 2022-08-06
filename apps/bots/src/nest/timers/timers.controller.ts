import { Controller } from '@nestjs/common';
import { type ClientProxyEventPayload, EventPattern } from '@tsuwari/shared';

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
  @EventPattern('streams.offline')
  streamOnline(data: ClientProxyEventPayload<'streams.online'> | ClientProxyEventPayload<'streams.offline'>) {
    this.service.handleStreamChange(data.channelId);
  }
}
