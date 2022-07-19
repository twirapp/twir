import { Controller } from '@nestjs/common';
import { EventPattern, MessagePattern, Payload } from '@nestjs/microservices';
import { ClientProxyCommands, ClientProxyCommandsKey, ClientProxyEvents, ClientProxyEventsKey } from '@tsuwari/shared';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  constructor(private readonly service: AppService) { }

  @EventPattern<ClientProxyEventsKey>('dota.cacheAccountsMatches')
  cacheAccountsMatches(@Payload() data: ClientProxyEvents['dota.cacheAccountsMatches']['input']) {
    this.service.getPresences(data);
  }

  @MessagePattern<ClientProxyCommandsKey>('dota.getProfileCard')
  getProfileCard(@Payload() data: ClientProxyCommands['dota.getProfileCard']['input']) {
    return this.service.getDotaProfileCard(data);
  }
}