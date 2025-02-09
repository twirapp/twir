import { ChannelCredentials, createChannel, createClient } from 'nice-grpc'

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper'
import { PORTS } from '../constants/constants.js'
import { WebsocketDefinition } from '../websockets/websockets'

import type { WebsocketClient } from '../websockets/websockets'

export async function createWebsocket(env: string): Promise<WebsocketClient> {
	const channel = createChannel(
		createClientAddr(env, 'websockets', PORTS.WEBSOCKET_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	)

	await waitReady(channel)

	const client = createClient(WebsocketDefinition, channel)

	return client as any
}
