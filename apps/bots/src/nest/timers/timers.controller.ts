import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { Observable, of } from 'rxjs';

import { addTimerToQueue, removeTimerFromQueue } from '../../libs/timers.js';

@Controller()
export class TimersController implements Bots.Timers {
  @GrpcMethod('Timers', 'addTimerToQueue')
  addTimerToQueue(data: Bots.AddTimerToQueueRequest): Observable<Bots.MockedResult> {
    try {
      console.log('controller', data.timerId);
      addTimerToQueue(data.timerId!);
      return of({ success: true });
    } catch (error) {
      return of({ success: false, error: (error as any).message });
    }
  }

  @GrpcMethod('Timers', 'removeTimerFromQueue')
  removeTimerFromQueue(data: Bots.RemoveTimerFromQueueRequest): Observable<Bots.MockedResult> {
    try {
      removeTimerFromQueue(data.timerId!);
      return of({ success: true });
    } catch (error) {
      return of({ success: false, error: (error as any).message });
    }
  }
}
