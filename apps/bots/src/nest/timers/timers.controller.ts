import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { Observable, of } from 'rxjs';

import { addTimerToQueue, removeTimerFromQueue } from '../../libs/timers.js';

@Controller()
export class TimersController implements Bots.Timers {
  @GrpcMethod('BotsServiceTimers', 'addTimersToQueue')
  addTimersToQueue(data: Bots.AddTimerToQueueRequest): Observable<Bots.MockedResult> {
    try {
      addTimerToQueue(data.timerId!);
      return of({ success: true });
    } catch (error) {
      return of({ success: false, error: (error as any).message });
    }
  }

  @GrpcMethod('BotsServiceTimers', 'addTimersToQueue')
  removeTimerFromQueue(data: Bots.RemoveTimerFromQueueRequest): Observable<Bots.MockedResult> {
    try {
      removeTimerFromQueue(data.timerId!);
      return of({ success: true });
    } catch (error) {
      return of({ success: false, error: (error as any).message });
    }
  }
}
