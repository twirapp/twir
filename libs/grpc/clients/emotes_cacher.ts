import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { PORTS } from '../constants/constants.js';
import {
	EmotesCacherClient,
	EmotesCacherDefinition,
} from '../emotes_cacher/emotes_cacher.client.js';

export const createEmotesCacher = async (env: string): Promise<EmotesCacherClient> => {
	const channel = createChannel(
		createClientAddr(env, 'emotes-cacher', PORTS.EMOTES_CACHER_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	);

	await waitReady(channel);

	const client = createClient(EmotesCacherDefinition, channel);

	return client as any;
};
