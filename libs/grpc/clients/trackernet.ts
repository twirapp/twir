import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { TrackernetClient, TrackernetDefinition } from '../generated/trackernet/trackernet.js';
import { PORTS } from '../servers/constants.js';
import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';

export const createTrackernet = async (env: string): Promise<TrackernetClient> => {
	const channel = createChannel(
		createClientAddr(env, 'trackernet', PORTS.TRACKERNET_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	);

	await waitReady(channel);

	const client = createClient(TrackernetDefinition, channel);

	return client as any;
};
