import { Controller } from '@nestjs/common';
import { Payload, EventPattern } from '@nestjs/microservices';

import { EventSub } from './eventsub.service.js';

@Controller()
export class EventSubController {
  constructor(private readonly service: EventSub) { }

  @EventPattern('eventsub.subscribeToEventsByChannelId')
  async subscribeToEvents(@Payload() channelId: string) {
    this.service.subscribeToEvents(channelId);
  }
}