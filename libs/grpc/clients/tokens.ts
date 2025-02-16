import { ChannelCredentials, createChannel, createClient } from 'nice-grpc'

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper'
import { PORTS } from '../constants/constants.js'
import { TokensDefinition } from '../tokens/tokens'

import type { TokensClient } from '../tokens/tokens'

export async function createTokens(env: string): Promise<TokensClient> {
	const channel = createChannel(
		createClientAddr(env, 'tokens', PORTS.TOKENS_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	)

	await waitReady(channel)

	const client = createClient(TokensDefinition, channel)

	return client as any
}
