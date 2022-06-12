import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Bots } from '@tsuwari/grpc';
import { of } from 'rxjs';

import * as DefCommands from '../../defaultCommands/index.js';

const commands = Object.values(DefCommands).flat();

@Controller('commands')
export class CommandsController implements Bots.Commands {
  @GrpcMethod('Commands', 'getDefaultCommands')
  getDefaultCommands() {
    return of({ commands: commands.map(c => ({ name: c.name, permission: c.permission })) });
  }
}
