import { ChannelCredentials, createChannel, createClient } from 'nice-grpc'

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper'
import { PORTS } from '../constants/constants.js'
import { YtsrDefinition } from '../ytsr/ytsr'

import type { YtsrClient } from '../ytsr/ytsr'

export async function createYtsr(env: string): Promise<YtsrClient> {
	const channel = createChannel(
		createClientAddr(env, 'ytsr', PORTS.YTSR_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	)

	await waitReady(channel)

	const client = createClient(YtsrDefinition, channel)

	return client as any
}
