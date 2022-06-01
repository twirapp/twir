import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots as GBots } from '@tsuwari/grpc';
import { Observable, of } from 'rxjs';

import { Bots } from '../../bots.js';

@Controller()
export class CommandsController implements GBots.Commands {
  @GrpcMethod('Commands', 'updateByChannelId')
  updateByChannelId(data: GBots.UpdateByChannelIdRequest): Observable<GBots.MockedResult> {
    try {
      Bots.updateCommandsCacheByChannelid(data.userId!);
      return of({ success: true });
    } catch (e) {
      return of({ success: false, error: (e as any).message });
    }
  }
}
