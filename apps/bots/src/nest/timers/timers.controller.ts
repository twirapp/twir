import { Controller } from '@nestjs/common';
import { EventPattern } from '@nestjs/microservices';

import { addTimerToQueue, removeTimerFromQueue } from '../../libs/timers.js';

@Controller()
export class TimersController {

  @EventPattern('bots.addTimerToQueue')
  addTimerToQueue(timerId: string) {
    addTimerToQueue(timerId);
  }

  @EventPattern('bots.removeTimerFromQueue')
  removeTimerFromQueue(timerId: string) {
    removeTimerFromQueue(timerId);
  }
}
