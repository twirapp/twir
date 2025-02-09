import { ChannelCredentials, createChannel, createClient } from 'nice-grpc'

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper'
import { PORTS } from '../constants/constants.js'
import { ParserDefinition } from '../parser/parser'

import type { ParserClient } from '../parser/parser'

export async function createParser(env: string): Promise<ParserClient> {
	const channel = createChannel(
		createClientAddr(env, 'parser', PORTS.PARSER_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	)

	await waitReady(channel)

	const client = createClient(ParserDefinition, channel)

	return client as any
}
