import io from 'socket.io-client'

import { claimDonation as claimDonationOnce } from '../libs/donation-dedupe.ts'
import type { ProviderTokens } from '../libs/provider-token-store.ts'
import type { Donate } from '../utils/onDonation.ts'

const SOCKET_URL = 'https://realtime.streamelements.com'
const MAX_RECONNECT_DELAY_MS = 30_000

export interface StreamElementsClientLike {
	readonly tokens: ProviderTokens
	refresh(): Promise<void>
}

export interface StreamElementsSocket {
	on(event: string, handler: (payload?: unknown) => void): StreamElementsSocket
	emit(event: string, payload: unknown): StreamElementsSocket
	close(): void
}

export interface StreamElementsSocketOptions {
	readonly transports: readonly string[]
	readonly reconnection: false
	readonly forceNew: true
}

export type StreamElementsSocketFactory = (
	url: string,
	options: StreamElementsSocketOptions
) => StreamElementsSocket

export interface StreamElementsConnectionOptions {
	readonly channelID: string
	readonly client: StreamElementsClientLike
	readonly socketFactory?: StreamElementsSocketFactory
	readonly onDonation?: (donation: Donate) => Promise<void>
	readonly claimDonation?: (provider: string, eventID: string) => Promise<boolean>
	readonly schedule?: (callback: () => void, milliseconds: number) => unknown
	readonly clearSchedule?: (handle: unknown) => void
	readonly random?: () => number
	readonly onError?: (error: unknown) => void
}

function defaultSocketFactory(
	url: string,
	options: StreamElementsSocketOptions
): StreamElementsSocket {
	return io.connect(url, options as SocketIOClient.ConnectOpts) as unknown as StreamElementsSocket
}

async function publishDonation(donation: Donate): Promise<void> {
	const { onDonation } = await import('../utils/onDonation.ts')
	await onDonation(donation)
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return value !== null && typeof value === 'object'
}

function optionalString(value: unknown): string | null {
	return typeof value === 'string' && value.length > 0 ? value : null
}

interface NormalizedTip {
	readonly eventID: string
	readonly amount: number | string
	readonly currency: string
	readonly message: string | null
	readonly userName: string | null
}

function normalizeTip(payload: unknown): NormalizedTip | null {
	if (!isRecord(payload) || payload.type !== 'tip' || !isRecord(payload.data)) return null
	const data = payload.data
	if ((typeof data.amount !== 'number' && typeof data.amount !== 'string')
		|| typeof data.currency !== 'string') {
		return null
	}
	const rawID = data.tipId ?? data._id ?? data.id
	return {
		eventID: typeof rawID === 'string' ? rawID : '',
		amount: data.amount,
		currency: data.currency,
		message: optionalString(data.message),
		userName: optionalString(data.username),
	}
}

export class StreamElementsConnection {
	readonly #channelID: string
	readonly #client: StreamElementsClientLike
	readonly #socketFactory: StreamElementsSocketFactory
	readonly #onDonation: (donation: Donate) => Promise<void>
	readonly #claimDonation: (provider: string, eventID: string) => Promise<boolean>
	readonly #schedule: (callback: () => void, milliseconds: number) => unknown
	readonly #clearSchedule: (handle: unknown) => void
	readonly #random: () => number
	readonly #onError: (error: unknown) => void
	#socket: StreamElementsSocket | null = null
	#reconnectHandle: unknown
	#generation = 0
	#reconnectAttempt = 0
	#refreshUsed = false
	#destroyed = false

	constructor(options: StreamElementsConnectionOptions) {
		this.#channelID = options.channelID
		this.#client = options.client
		this.#socketFactory = options.socketFactory ?? defaultSocketFactory
		this.#onDonation = options.onDonation ?? publishDonation
		this.#claimDonation = options.claimDonation ?? claimDonationOnce
		this.#schedule = options.schedule ?? ((callback, milliseconds) => setTimeout(callback, milliseconds))
		this.#clearSchedule = options.clearSchedule ?? ((handle) => {
			clearTimeout(handle as ReturnType<typeof setTimeout>)
		})
		this.#random = options.random ?? Math.random
		this.#onError = options.onError ?? console.error
	}

	connect(): void {
		if (this.#destroyed) return
		this.#replaceSocket()
	}

	destroy(): void {
		this.#destroyed = true
		this.#generation += 1
		this.#cancelReconnect()
		this.#socket?.close()
		this.#socket = null
	}

	#replaceSocket(): void {
		if (this.#destroyed) return
		this.#generation += 1
		const generation = this.#generation
		this.#socket?.close()
		const socket = this.#socketFactory(SOCKET_URL, {
			transports: ['websocket'],
			reconnection: false,
			forceNew: true,
		})
		this.#socket = socket

		socket.on('connect', () => {
			if (!this.#isCurrent(generation, socket)) return
			this.#reconnectAttempt = 0
			socket.emit('authenticate', {
				method: 'oauth2',
				token: this.#client.tokens.accessToken,
			})
		})
		socket.on('event', (payload) => {
			if (!this.#isCurrent(generation, socket)) return
			void this.#handleEvent(payload).catch(this.#onError)
		})
		socket.on('unauthorized', () => {
			if (!this.#isCurrent(generation, socket)) return
			void this.#handleUnauthorized(socket).catch(this.#onError)
		})
		const reconnect = () => {
			if (!this.#isCurrent(generation, socket)) return
			this.#scheduleReconnect(generation)
		}
		socket.on('disconnect', reconnect)
		socket.on('connect_error', reconnect)
	}

	#isCurrent(generation: number, socket: StreamElementsSocket): boolean {
		return !this.#destroyed && generation === this.#generation && socket === this.#socket
	}

	async #handleEvent(payload: unknown): Promise<void> {
		const tip = normalizeTip(payload)
		if (!tip) return
		if (!await this.#claimDonation('streamelements', tip.eventID)) return
		await this.#onDonation({
			twitchUserId: this.#channelID,
			amount: tip.amount,
			currency: tip.currency,
			message: tip.message,
			userName: tip.userName,
		})
	}

	async #handleUnauthorized(socket: StreamElementsSocket): Promise<void> {
		this.#cancelReconnect()
		if (this.#refreshUsed) {
			this.#generation += 1
			socket.close()
			this.#socket = null
			return
		}
		this.#refreshUsed = true
		this.#generation += 1
		socket.close()
		this.#socket = null
		await this.#client.refresh()
		if (!this.#destroyed) this.#replaceSocket()
	}

	#scheduleReconnect(generation: number): void {
		if (this.#reconnectHandle !== undefined) return
		const exponential = Math.min(1_000 * 2 ** this.#reconnectAttempt, MAX_RECONNECT_DELAY_MS)
		const jittered = Math.round(exponential * (0.5 + this.#random()))
		const delay = Math.min(jittered, MAX_RECONNECT_DELAY_MS)
		this.#reconnectAttempt += 1
		this.#reconnectHandle = this.#schedule(() => {
			this.#reconnectHandle = undefined
			if (this.#destroyed || generation !== this.#generation) return
			this.#replaceSocket()
		}, delay)
	}

	#cancelReconnect(): void {
		if (this.#reconnectHandle === undefined) return
		this.#clearSchedule(this.#reconnectHandle)
		this.#reconnectHandle = undefined
	}
}
