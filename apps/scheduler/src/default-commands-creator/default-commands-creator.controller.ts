import { Controller } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import { MessagePattern, type ClientProxyCommandPayload } from '@tsuwari/shared';

import { DefaultCommandsCreatorService } from './default-commands-creator.service.js';

@Controller()
export class DefaultCommandsCreatorController {
  constructor(private readonly service: DefaultCommandsCreatorService) { }

  @MessagePattern('bots.createDefaultCommands')
  async setCommandCache(@Payload() data: ClientProxyCommandPayload<'bots.createDefaultCommands'>) {
    // return await this.service.createDefaultCommands(data);
    return true
  }
}