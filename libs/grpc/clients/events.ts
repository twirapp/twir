import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { PORTS } from '../constants/constants.js';
import { EventsClient, EventsDefinition } from '../generated/events/events.js';

export const createEvents = async (env: string): Promise<EventsClient> => {
	const channel = createChannel(
		createClientAddr(env, 'events', PORTS.EVENTS_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	);

	await waitReady(channel);

	const client = createClient(EventsDefinition, channel);

	return client as any;
};
