import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { ParserClient, ParserDefinition } from '../generated/parser/parser.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createParser = async (env: string): Promise<ParserClient> => {
  const channel = createChannel(
    createClientAddr(env, 'parser', PORTS.PARSER_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  await waitReady(channel);

  const client = createClient(ParserDefinition, channel);

  return client as any;
};
