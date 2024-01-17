import { ChannelCredentials, createChannel, createClient } from 'nice-grpc';

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js';
import { PORTS } from '../constants/constants.js';
import { TokensClient, TokensDefinition } from '../dist/tokens/tokens.js';

export const createTokens = async (env: string): Promise<TokensClient> => {
	const channel = createChannel(
		createClientAddr(env, 'tokens', PORTS.TOKENS_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	);

	await waitReady(channel);

	const client = createClient(TokensDefinition, channel);

	return client as any;
};
