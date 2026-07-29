import { newBus } from '@twir/bus-core'
import { config } from '@twir/config'
import { connect } from 'nats'
import { Innertube } from 'youtubei.js'

import { closeDatabase, ensureYoutubeChatter, getYoutubeBinding, listYoutubeBindings } from './db.ts'
import { LiveChatManager } from './live-chat.ts'
import { RedisBindingOwnership } from './locks.ts'

import type { LiveChatSource } from './live-chat.ts'

const RECONCILE_INTERVAL_MS = Number.parseInt(Bun.env.YTSUB_RECONCILE_INTERVAL_MS ?? '120000', 10)
const NATS_URL = Bun.env.NODE_ENV === 'production' ? 'nats://nats:4222' : 'nats://127.0.0.1:4222'

const yt = await Innertube.create()
const nc = await connect({ servers: NATS_URL })
const bus = newBus(nc)
const liveChatSource: LiveChatSource = {
	async resolve(binding) {
		const endpoint = await yt.resolveURL(`https://www.youtube.com/channel/${binding.platformChannelId}/live`)
		const info = await yt.getInfo(endpoint)
		const liveChat = info.getLiveChat()
		return {
			broadcasterName: info.basic_info.channel?.name ?? info.basic_info.author ?? binding.platformChannelId,
			session: {
				onStart(listener): void {
					liveChat.on('start', listener)
				},
				onChatUpdate(listener): void {
					liveChat.on('chat-update', listener)
				},
				onError(listener): void {
					liveChat.on('error', listener)
				},
				onEnd(listener): void {
					liveChat.on('end', listener)
				},
				selectLiveChat(): void {
					liveChat.applyFilter('LIVE_CHAT')
				},
				start(): void {
					liveChat.start()
				},
				stop(): void {
					liveChat.stop()
				},
			},
		}
	},
}
const ownership = new RedisBindingOwnership(new Bun.RedisClient(config.REDIS_URL))
const liveChats = new LiveChatManager(liveChatSource, bus, {
	ensureChatter: ensureYoutubeChatter,
	ownership,
})

async function subscribeChannel(channelId: string): Promise<void> {
	const binding = await getYoutubeBinding(channelId)
	if (!binding || !binding.enabled) {
		return
	}
	await liveChats.subscribe(binding)
}

await bus.EventSub.SubscribeToAllEvents.subscribeGroup('ytsub', async (request) => {
	if (request.Platform === '' || request.Platform === 'youtube') {
		await subscribeChannel(request.ChannelID)
	}
	return {}
})

await bus.EventSub.Unsubscribe.subscribeGroup('ytsub', async (request) => {
	if (request.Platform !== '' && request.Platform !== 'youtube') {
		return {}
	}
	if (request.Binding) {
		await liveChats.unsubscribe(request.Binding.ID)
		return {}
	}
	const binding = await getYoutubeBinding(request.ChannelID)
	if (binding) {
		await liveChats.unsubscribe(binding.id)
	}
	return {}
})

async function reconcile(): Promise<void> {
	const bindings = await listYoutubeBindings()
	await liveChats.reconcile(bindings)
}

await reconcile()

const reconcileTimer = setInterval(() => {
	void reconcile().catch((error: unknown) => {
		console.error('youtube.reconcile.failed', { error })
	})
}, RECONCILE_INTERVAL_MS)

let shutdownPromise: Promise<void> | undefined

function shutdown(): Promise<void> {
	if (shutdownPromise) {
		return shutdownPromise
	}
	shutdownPromise = (async (): Promise<void> => {
		clearInterval(reconcileTimer)
		try {
			await liveChats.close()
			await nc.drain()
		} catch (error) {
			const drainError = error instanceof Error ? error : new Error('NATS drain failed')
			console.error('youtube.shutdown.drain.failed', { error: drainError })
		} finally {
			try {
				await closeDatabase()
			} catch (error) {
				const closeError = error instanceof Error ? error : new Error('Database close failed')
				console.error('youtube.shutdown.database-close.failed', { error: closeError })
			}
		}
	})()
	return shutdownPromise
}

process.once('SIGINT', () => {
	void shutdown().catch((error: unknown) => {
		console.error('youtube.shutdown.failed', { error })
	})
})
process.once('SIGTERM', () => {
	void shutdown().catch((error: unknown) => {
		console.error('youtube.shutdown.failed', { error })
	})
})

console.info('YouTube live chat subscriber started')
