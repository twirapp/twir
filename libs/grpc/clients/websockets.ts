import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { WebsocketClient, WebsocketDefinition } from '../generated/websockets/websockets.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createWebsocket = async (env: string): Promise<WebsocketClient> => {
  const channel = createChannel(
    createClientAddr(env, 'websockets', PORTS.WEBSOCKET_SERVER_PORT),
    ChannelCredentials.createInsecure(),
    CLIENT_OPTIONS,
  );

  await waitReady(channel);

  const client = createClient(WebsocketDefinition, channel);

  return client as any;
};
