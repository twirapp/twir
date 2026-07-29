import { newBus } from '@twir/bus-core'
import { connect } from 'nats'
import { Innertube } from 'youtubei.js'

import { closeDatabase, getYoutubeBinding, listYoutubeBindings } from './db.ts'
import { LiveChatManager } from './live-chat.ts'

const RECONCILE_INTERVAL_MS = Number.parseInt(Bun.env.YTSUB_RECONCILE_INTERVAL_MS ?? '120000', 10)
const NATS_URL = Bun.env.NODE_ENV === 'production' ? 'nats://nats:4222' : 'nats://127.0.0.1:4222'

const yt = await Innertube.create()
const nc = await connect({ servers: NATS_URL })
const bus = newBus(nc)
const liveChats = new LiveChatManager(yt, bus)

async function subscribeChannel(channelId: string): Promise<void> {
	const binding = await getYoutubeBinding(channelId)
	if (!binding || !binding.enabled) {
		return
	}
	await liveChats.subscribe(binding)
}

await Promise.all((await listYoutubeBindings()).map((binding) => liveChats.subscribe(binding)))

void bus.EventSub.SubscribeToAllEvents.subscribeGroup('ytsub', async (request) => {
	if (request.Platform === '' || request.Platform === 'youtube') {
		await subscribeChannel(request.ChannelID)
	}
	return {}
})

void bus.EventSub.Unsubscribe.subscribeGroup('ytsub', async (request) => {
	if (request.Platform !== '' && request.Platform !== 'youtube') {
		return {}
	}
	if (request.Binding) {
		liveChats.unsubscribe(request.Binding.ID)
		return {}
	}
	const binding = await getYoutubeBinding(request.ChannelID)
	if (binding) {
		liveChats.unsubscribe(binding.id)
	}
	return {}
})

const reconcile = (): void => {
	void liveChats.reconcile()
}

const reconcileTimer = setInterval(reconcile, RECONCILE_INTERVAL_MS)

async function shutdown(): Promise<void> {
	clearInterval(reconcileTimer)
	liveChats.close()
	await nc.drain()
	await closeDatabase()
}

process.once('SIGINT', () => void shutdown())
process.once('SIGTERM', () => void shutdown())

console.info('YouTube live chat subscriber started')
