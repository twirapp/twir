import { sleep } from 'bun'
import Centrifuge from 'centrifuge'
import { RateLimiter, RedisStore } from 'rate-limiter-algorithms'
import WebSocket from 'ws'

import { client } from '../libs/redis'
import { onDonation } from '../utils/onDonation'

import type { Subscription } from 'centrifuge'

export const globalRequestLimiter = new RateLimiter({
	store: new RedisStore({
		prefix: 'integrations:rla:',
		rawCall: (...args) => client.sendCommand(args),
	}),
	algorithm: 'sliding-window-counter',
	limit: 50,
	windowMs: 1 * 60 * 1000,
})

export class DonationAlerts {
	#socket: Centrifuge | null
	#channel: Subscription | null

	constructor(
		private readonly accessToken: string,
		private readonly donationAlertsUserId: string,
		private readonly socketConnectionToken: string,
		private readonly twitchUserId: string,
	) {}

	async init() {
		this.#socket = new Centrifuge('wss://centrifugo.donationalerts.com/connection/websocket', {
			websocket: WebSocket,
			onPrivateSubscribe: async (ctx, cb) => {
				while (true) {
					const { isAllowed } = await globalRequestLimiter.consume(this.twitchUserId)
					if (!isAllowed) {
						await sleep(1000)
						continue
					}

					const request = await fetch('https://www.donationalerts.com/api/v1/centrifuge/subscribe', {
						method: 'POST',
						body: JSON.stringify(ctx.data),
						headers: { Authorization: `Bearer ${this.accessToken}` },
					})

					const response = await request.json()
					if (!request.ok) {
						console.error(response)
						// eslint-disable-next-line ts/ban-ts-comment
						// @ts-expect-error
						cb({ status: request.status, data: {} })
						break
					}

					cb({ status: 200, data: { channels: response.channels } })
					break
				}
			},
		})

		this.#socket.setToken(this.socketConnectionToken)
		this.#socket.connect()

		this.#socket.on('connect', () => {
			console.info(`Connected to donationAlerts #${this.donationAlertsUserId}`)
		})

		this.#channel = this.#socket.subscribe(`$alerts:donation_${this.donationAlertsUserId}`)

		this.#channel.on('publish', async ({ data }) => this.#donateCallback(data))

		return this
	}

	async #donateCallback(data: DonationAlertsMessage) {
		console.info(`[DONATIONALERTS #${this.twitchUserId}]  Donation from ${data.username}: ${data.amount} ${data.currency}`)
		await onDonation({
			twitchUserId: this.twitchUserId,
			amount: data.amount,
			currency: data.currency,
			message: data.message,
			userName: data.username,
		})
	}

	async destroy() {
		this.#channel?.removeAllListeners()?.unsubscribe()
		this.#socket?.removeAllListeners()?.disconnect()

		this.#socket = null
		this.#channel = null
	}
}

export interface DonationAlertsMessage {
	id: number
	name: string
	username?: string | null
	message: string | null
	message_type: 'text' | 'audio'
	payin_system: null | any
	amount: number
	currency: string
	amount_in_user_currency: number
	recipient_name: string
	recipient: {
		user_id: number
		code: string
		name: string
		avatar: string
	}
	created_at: string
	shown_at: null | any
	reason: string
}
