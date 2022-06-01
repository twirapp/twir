import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots as GBots } from '@tsuwari/grpc';
import { Observable, of } from 'rxjs';

import { Bots } from '../../bots.js';

@Controller()
export class GreetingsController implements GBots.Greetings {
  @GrpcMethod('BotsServiceGreetings', 'updateByChannelId')
  updateByChannelId(request: GBots.UpdateByChannelIdRequest): Observable<GBots.MockedResult> {
    try {
      Bots.updateGreetingsCacheByChannelid(request.userId!);
      return of({ success: true });
    } catch (error) {
      return of({ success: false, error: (error as any).message });
    }
  }
}
