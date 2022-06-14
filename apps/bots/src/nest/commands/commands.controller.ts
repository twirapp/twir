import { Controller } from '@nestjs/common';
import { MessagePattern } from '@nestjs/microservices';
import { ClientProxyResult } from '@tsuwari/shared';
import { of } from 'rxjs';

import * as DefCommands from '../../defaultCommands/index.js';

const commands = Object.values(DefCommands).flat();

@Controller()
export class CommandsController {
  @MessagePattern('bots.getDefaultCommands')
  getDefaultCommands(): ClientProxyResult<'bots.getDefaultCommands'> {
    return of(commands.map(c => ({ name: c.name, permission: c.permission, description: c.description })));
  }
}
