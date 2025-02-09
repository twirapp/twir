import { ChannelCredentials, createChannel, createClient } from 'nice-grpc'

import { CLIENT_OPTIONS, createClientAddr, waitReady } from './helper.js'
import { PORTS } from '../constants/constants.js'
import { EventsDefinition } from '../events/events'

import type { EventsClient } from '../events/events'

export async function createEvents(env: string): Promise<EventsClient> {
	const channel = createChannel(
		createClientAddr(env, 'events', PORTS.EVENTS_SERVER_PORT),
		ChannelCredentials.createInsecure(),
		CLIENT_OPTIONS,
	)

	await waitReady(channel)

	const client = createClient(EventsDefinition, channel)

	return client as any
}
