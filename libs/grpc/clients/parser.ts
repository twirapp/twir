import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { ParserClient, ParserDefinition } from '../generated/parser/parser.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr } from './helper.js';

export const createParser = (env: string): ParserClient => {
  const channel = createChannel(
    createClientAddr(env, 'parser', PORTS.PARSER_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  const client = createClient(ParserDefinition, channel);

  return client;
};
