import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { PORTS } from '../constants/constants.js';
import { ParserClient, ParserDefinition } from '../dist/parser/parser.js';

export const createParser = async (env: string): Promise<ParserClient> => {
	const channel = createChannel(
		createClientAddr(env, 'parser', PORTS.PARSER_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	);

	await waitReady(channel);

	const client = createClient(ParserDefinition, channel);

	return client as any;
};
