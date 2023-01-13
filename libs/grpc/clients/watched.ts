import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { WatchedClient, WatchedDefinition } from '../generated/watched/watched.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createWatched = async (env: string): Promise<WatchedClient> => {
  const channel = createChannel(
    createClientAddr(env, 'watched', PORTS.WATCHED_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  await waitReady(channel);

  const client = createClient(WatchedDefinition, channel);

  return client as any;
};
