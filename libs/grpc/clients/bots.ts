import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { BotsClient, BotsDefinition } from '../bots/bots.client.js';
import { PORTS } from '../constants/constants.js';

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
