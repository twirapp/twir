import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots as GBots } from '@tsuwari/grpc';
import { Observable, of } from 'rxjs';

import { Bots } from '../bots.js';
import * as Variables from '../parser/modules/index.js';

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

  @GrpcMethod('Main', 'GetVariables')
  getVariables(): Observable<GBots.GetVariablesResult> {
    const variables = Object.values(Variables).map(v => {
      const modules = Array.isArray(v) ? v : [v];

      return modules
        .filter(m => typeof m.visible !== 'undefined' ? m.visible : true)
        .map(m => ({ name: m.key, example: m.example, description: m.description }));
    }).flat();

    return of({ variables });
  }
}
