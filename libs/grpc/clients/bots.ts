import { ChannelCredentials, createChannel, createClient, waitForChannelReady } from 'nice-grpc';

import { BotsClient, BotsDefinition } from '../generated/bots/bots.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createBots = async (env: string): Promise<BotsClient> => {
  const channel = createChannel(
    createClientAddr(env, 'bots', PORTS.BOTS_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  await waitReady(channel);

  const client = createClient(BotsDefinition, channel);

  return client as any;
};
