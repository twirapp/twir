import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { WatchedClient, WatchedDefinition } from '../generated/watched/watched.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr } from './helper.js';

export const createWatched = (env: string): WatchedClient => {
  const channel = createChannel(
    createClientAddr(env, 'watched', PORTS.WATCHED_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  const client = createClient(WatchedDefinition, channel);

  return client as any;
};
