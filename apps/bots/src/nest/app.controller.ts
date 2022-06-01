import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots as GBots } from '@tsuwari/grpc';
import { Observable, of } from 'rxjs';

import { Bots } from '../bots.js';

@Controller()
export class AppController implements GBots.Main {
  @GrpcMethod('Main', 'joinOrLeave')
  joinOrLeave(request: GBots.JoinOrLeaveRequest): Observable<GBots.MockedResult> {
    const bot = Bots.cache.get(request.botId!);
    if (bot) {
      const action = request.action === GBots.JoinOrLeaveRequest.Action.JOIN ? 'join' : 'part';
      bot[action](request.username!);
      return of({ success: true });
    } else {
      return of({ success: false, error: 'Bot not found for that user.' });
    }
  }
}
