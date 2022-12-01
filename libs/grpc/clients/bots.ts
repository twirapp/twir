import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { BotsClient, BotsDefinition } from '../generated/bots/bots.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr } from './helper.js';

export const createBots = (env: string): BotsClient => {
  const channel = createChannel(
    createClientAddr(env, 'bots', PORTS.BOTS_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  const client = createClient(BotsDefinition, channel);

  return client;
};
