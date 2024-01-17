import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { PORTS } from '../constants/constants.js';
import { WebsocketClient, WebsocketDefinition } from '../websockets/websockets.client.js';

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
