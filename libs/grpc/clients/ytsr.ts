import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { PORTS } from '../constants/constants.js';
import { YtsrClient, YtsrDefinition } from '../generated/ytsr/ytsr.js';

export const createYtsr = async (env: string): Promise<YtsrClient> => {
	const channel = createChannel(
		createClientAddr(env, 'ytsr', PORTS.YTSR_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	);

	await waitReady(channel);

	const client = createClient(YtsrDefinition, channel);

	return client as any;
};
