import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { TokensClient, TokensDefinition } from '../generated/tokens/tokens.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createTokens = async (env: string): Promise<TokensClient> => {
  const channel = createChannel(
    createClientAddr(env, 'tokens', PORTS.TOKENS_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  await waitReady(channel);

  const client = createClient(TokensDefinition, channel);

  return client as any;
};
