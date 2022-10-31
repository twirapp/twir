import * as EventSub from '@tsuwari/nats/eventsub';
import { NatsConnection } from 'nats';

import { EventSub as EventSubService } from '../eventsub/eventsub.service.js';
import { app } from '../index.js';

export async function listenForDefaultCommands(nats: NatsConnection) {
  for await (const event of nats.subscribe(EventSub.SUBJECTS.SUBSCTUBE_TO_EVENTS_BY_CHANNEL_ID)) {
    const data = EventSub.SubscribeToEvents.fromBinary(event.data);
    const service = app.get(EventSubService);
    await service.subscribeToEvents(data.channelId);
    event.respond();
  }
}
