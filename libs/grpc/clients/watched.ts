import { ChannelCredentials, createChannel, createClient, RawClient } from 'nice-grpc';
import { FromTsProtoServiceDefinition } from 'nice-grpc/lib/service-definitions/ts-proto.js';

import { WatchedDefinition } from '../generated/watched/watched.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr } from './helper.js';

export const createWatched = (
  env: string,
): RawClient<FromTsProtoServiceDefinition<WatchedDefinition>> => {
  const channel = createChannel(
    createClientAddr(env, 'watched', PORTS.WATCHED_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  const client = createClient(WatchedDefinition, channel);

  return client;
};
