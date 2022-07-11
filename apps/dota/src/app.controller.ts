import { Controller } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import { ClientProxyEvents, ClientProxyEventsKey } from '@tsuwari/shared';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  constructor(private readonly service: AppService) { }

  @EventPattern<ClientProxyEventsKey>('dota.cacheAccountsMatches')
  cacheAccountsMatches(@Payload() data: ClientProxyEvents['dota.cacheAccountsMatches']['input']) {
    this.service.getPresences(data);
  }
}