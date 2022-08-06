import { Controller } from '@nestjs/common';
import { Payload } from '@nestjs/microservices';
import { type ClientProxyEventPayload, EventPattern } from '@tsuwari/shared';

import { EventSub } from './eventsub.service.js';

@Controller()
export class EventSubController {
  constructor(private readonly service: EventSub) { }

  @EventPattern('eventsub.subscribeToEventsByChannelId')
  async subscribeToEvents(@Payload() channelId: ClientProxyEventPayload<'eventsub.subscribeToEventsByChannelId'>) {
    this.service.subscribeToEvents(channelId);
  }
}