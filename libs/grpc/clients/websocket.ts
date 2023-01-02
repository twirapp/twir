import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { WebsocketClient, WebsocketDefinition } from '../generated/websocket/websocket.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr } from './helper.js';

export const createWebsocket = (env: string): WebsocketClient => {
  const channel = createChannel(
    createClientAddr(env, 'websocket', PORTS.WEBSOCKET_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  const client = createClient(WebsocketDefinition, channel);

  return client as any;
};
