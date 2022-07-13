import { Controller } from '@nestjs/common';
import { Payload, MessagePattern } from '@nestjs/microservices';
import { ClientProxyCommands, ClientProxyCommandsKey } from '@tsuwari/shared';

import { DefaultCommandsCreatorService } from './default-commands-creator.service.js';

@Controller()
export class DefaultCommandsCreatorController {
  constructor(private readonly service: DefaultCommandsCreatorService) { }

  @MessagePattern<ClientProxyCommandsKey>('bots.createDefaultCommands')
  async setCommandCache(@Payload() data: ClientProxyCommands['bots.createDefaultCommands']['input']) {
    return await this.service.createDefaultCommands(data);
  }
}