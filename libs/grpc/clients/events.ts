import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { EventsClient, EventsDefinition } from '../generated/events/events.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createEvents = async (env: string): Promise<EventsClient> => {
  const channel = createChannel(
    createClientAddr(env, 'events', PORTS.EVENTS_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  await waitReady(channel);

  const client = createClient(EventsDefinition, channel);

  return client as any;
};
