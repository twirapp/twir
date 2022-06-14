import { Controller } from '@nestjs/common';
import { EventPattern, MessagePattern, Payload } from '@nestjs/microservices';
import { ClientProxyResult } from '@tsuwari/shared';
import { of } from 'rxjs';

import { Bots } from '../bots.js';
import * as Variables from '../parser/modules/index.js';

@Controller()
export class AppController {
  @EventPattern('bots.joinOrLeave')
  joinOrLeave(@Payload() data: { action: 'join' | 'part', botId: string, username: string }) {
    const bot = Bots.cache.get(data.botId);
    if (bot) {
      bot[data.action](data.username);
    }
  }

  @MessagePattern('bots.getVariables')
  getVariables(): ClientProxyResult<'bots.getVariables'> {
    const variables = Object.values(Variables).map(v => {
      const modules = Array.isArray(v) ? v : [v];

      return modules
        .filter(m => typeof m.visible !== 'undefined' ? m.visible : true)
        .map(m => ({ name: m.key, example: m.example, description: m.description }));
    }).flat();

    return of(variables);
  }
}
