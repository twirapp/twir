import { Controller } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import {
  EventPattern,
  MessagePattern,
  type ClientProxyCommandPayload,
  type ClientProxyEventPayload
} from '@tsuwari/shared';

import { AppService } from './app.service.js';

@Controller()
export class AppController {
  constructor(private readonly service: AppService) {}

  @EventPattern('dota.cacheAccountsMatches')
  cacheAccountsMatches(@Payload() data: ClientProxyEventPayload<'dota.cacheAccountsMatches'>) {
    this.service.getPresences(data);
  }

  // @MessagePattern('dota.getProfileCard')
  // getProfileCard(@Payload() data: ClientProxyCommandPayload<'dota.getProfileCard'>) {
  //   return this.service.getDotaProfileCard(data);
  // }
}
