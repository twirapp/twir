import io from 'socket.io-client'

import { claimDonation as claimDonationOnce } from '../libs/donation-dedupe.ts'
import type { Donate } from '../utils/onDonation.ts'

const SOCKET_URL = 'https://sockets.streamlabs.com'
const MAX_RECONNECT_DELAY_MS = 30_000

export interface StreamLabsSocket {
	on(event: string, handler: (payload?: unknown) => void): StreamLabsSocket
	close(): void
}

export interface StreamLabsSocketOptions {
	readonly transports: readonly string[]
	readonly reconnection: false
	readonly forceNew: true
}

export type StreamLabsSocketFactory = (
	url: string,
	options: StreamLabsSocketOptions
) => StreamLabsSocket

export interface StreamLabsConnectionOptions {
	readonly channelID: string
	readonly socketToken: string
	readonly socketFactory?: StreamLabsSocketFactory
	readonly onDonation?: (donation: Donate) => Promise<void>
	readonly claimDonation?: (provider: string, eventID: string) => Promise<boolean>
	readonly schedule?: (callback: () => void, milliseconds: number) => unknown
	readonly clearSchedule?: (handle: unknown) => void
	readonly random?: () => number
	readonly onError?: (error: unknown) => void
}

interface NormalizedDonation {
	readonly eventID: string
	readonly amount: number | string
	readonly currency: string
	readonly message: string | null
	readonly userName: string | null
}

function defaultSocketFactory(url: string, options: StreamLabsSocketOptions): StreamLabsSocket {
	return io.connect(url, options as SocketIOClient.ConnectOpts) as unknown as StreamLabsSocket
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

function normalizeDonations(payload: unknown): NormalizedDonation[] {
	if (!isRecord(payload) || payload.type !== 'donation' || !Array.isArray(payload.message)) return []
	const eventID = optionalString(payload.event_id)
	const multipleMessages = payload.message.length > 1
	const donations: NormalizedDonation[] = []
	for (const [index, rawMessage] of payload.message.entries()) {
		if (!isRecord(rawMessage)
			|| (typeof rawMessage.amount !== 'number' && typeof rawMessage.amount !== 'string')
			|| typeof rawMessage.currency !== 'string') {
			continue
		}
		const messageID = optionalString(rawMessage._id)
		donations.push({
			eventID: messageID ?? (eventID ? `${eventID}${multipleMessages ? `:${index}` : ''}` : ''),
			amount: rawMessage.amount,
			currency: rawMessage.currency,
			message: optionalString(rawMessage.message),
			userName: optionalString(rawMessage.from),
		})
	}
	return donations
}

export class StreamLabsConnection {
	readonly #channelID: string
	readonly #socketToken: string
	readonly #socketFactory: StreamLabsSocketFactory
	readonly #onDonation: (donation: Donate) => Promise<void>
	readonly #claimDonation: (provider: string, eventID: string) => Promise<boolean>
	readonly #schedule: (callback: () => void, milliseconds: number) => unknown
	readonly #clearSchedule: (handle: unknown) => void
	readonly #random: () => number
	readonly #onError: (error: unknown) => void
	#socket: StreamLabsSocket | null = null
	#reconnectHandle: unknown
	#generation = 0
	#reconnectAttempt = 0
	#destroyed = false

	constructor(options: StreamLabsConnectionOptions) {
		this.#channelID = options.channelID
		this.#socketToken = options.socketToken
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
		const query = new URLSearchParams({ token: this.#socketToken })
		const socket = this.#socketFactory(`${SOCKET_URL}?${query}`, {
			transports: ['websocket'],
			reconnection: false,
			forceNew: true,
		})
		this.#socket = socket

		socket.on('connect', () => {
			if (!this.#isCurrent(generation, socket)) return
			this.#reconnectAttempt = 0
		})
		socket.on('event', (payload) => {
			if (!this.#isCurrent(generation, socket)) return
			void this.#handleEvent(payload, generation, socket).catch(this.#onError)
		})
		const reconnect = () => {
			if (!this.#isCurrent(generation, socket)) return
			this.#scheduleReconnect(generation)
		}
		socket.on('disconnect', reconnect)
		socket.on('connect_error', reconnect)
	}

	#isCurrent(generation: number, socket: StreamLabsSocket): boolean {
		return !this.#destroyed && generation === this.#generation && socket === this.#socket
	}

	async #handleEvent(
		payload: unknown,
		generation: number,
		socket: StreamLabsSocket
	): Promise<void> {
		const donations = normalizeDonations(payload)
		await Promise.all(donations.map(async (donation) => {
			try {
				if (!await this.#claimDonation('streamlabs', donation.eventID)) return
				if (!this.#isCurrent(generation, socket)) return
				await this.#onDonation({
					twitchUserId: this.#channelID,
					amount: donation.amount,
					currency: donation.currency,
					message: donation.message,
					userName: donation.userName,
				})
			} catch (error) {
				this.#onError(error)
			}
		}))
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
