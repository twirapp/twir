import { sleep } from 'bun'
import Centrifuge from 'centrifuge'
import { RateLimiter, RedisStore } from 'rate-limiter-algorithms'
import ws from 'ws'
// eslint-disable-next-line ts/ban-ts-comment
// @ts-expect-error
import { XMLHttpRequest } from 'xmlhttprequest'

import { client } from '../libs/redis.ts'
import { onDonation } from '../utils/onDonation.js'

import type { Subscription } from 'centrifuge'

// eslint-disable-next-line no-restricted-globals
global.XMLHttpRequest = XMLHttpRequest

const requestLimiter = new RateLimiter({
	store: new RedisStore({
		prefix: 'integrations:donatepay:rla:',
		rawCall: async (...args) => {
			if (!args.at(0)) return

			return await client.send(args.at(0)!, args.slice(1))
		},
	}),
	algorithm: 'sliding-window-counter',
	limit: 1,
	windowMs: 1000,
})

export class DonatePay {
	#centrifuge: Centrifuge
	#subscription: Subscription
	#timeout: Timer

	constructor(
		private readonly twitchUserId: string,
		private readonly apiKey: string,
		private baseDomain = 'donatepay.ru'
	) {}

	get isEu() {
		return this.baseDomain === 'donatepay.eu'
	}

	async connect() {
		if (this.#centrifuge || this.#subscription) {
			await this.disconnect()
		}

		const userData = await this.#getUserData()
		if (!userData) {
			console.error(`DonatePay: something went wrong when getting token of ${this.twitchUserId}`)
			return
		}

		this.#centrifuge = new Centrifuge(
			`wss://centrifugo.${this.baseDomain}:43002/connection/websocket`,
			{
				subscribeEndpoint: `https://${this.baseDomain}/api/v2/socket/token`,
				subscribeParams: {
					access_token: this.apiKey,
				},
				disableWithCredentials: true,
				websocket: ws,
				ping: true,
				pingInterval: 5000,
			}
		)

		this.#centrifuge.setToken(userData.token)

		this.#subscription = this.#centrifuge.subscribe(`$public:${userData.id}`, (data) => {
			return this.#eventCallback(data)
		})

		const logDisconnect = (args: any[]) =>
			console.info(`DonatePay(${this.twitchUserId}): disconnected`, args)

		this.#centrifuge.on('disconnect', logDisconnect)
		this.#subscription.on('disconnect', logDisconnect)

		this.#centrifuge.on('connect', () => {
			console.info(`DonatePay: connected to channel ${this.twitchUserId}`)
		})

		this.#centrifuge.connect()
		this.#timeout = setTimeout(() => this.connect(), 10 * 60 * 1000)
	}

	async #eventCallback(message: DonatePayEvent) {
		if (message.data.notification.type !== 'donation') return

		const { vars } = message.data.notification

		await onDonation({
			twitchUserId: this.twitchUserId,
			amount: vars.sum,
			currency: vars.currency,
			message: vars.comment,
			userName: vars.name,
		})
	}

	async disconnect() {
		clearTimeout(this.#timeout)
		this.#subscription?.unsubscribe()
		this.#centrifuge?.disconnect()
		console.info(`DonatePay: disconnected from channel ${this.twitchUserId}`)
	}

	async #getUserId() {
		while (true) {
			const { isAllowed } = await requestLimiter.consume('donatepay')
			if (!isAllowed) {
				await sleep(1000)
				continue
			}

			const getUserParams = new URLSearchParams({
				access_token: this.apiKey,
			})
			const req = await fetch(`https://${this.baseDomain}/api/v1/user?${getUserParams}`)
			if (req.status === 429) {
				await sleep(1000)
				continue
			}

			if (!req.ok) {
				console.error(`DonatePay: cannot get userId for ${this.twitchUserId}`, await req.text())
				throw new Error(`cannot get userId for ${this.twitchUserId}`)
			}

			const data = await req.json()
			if (!data.data?.id) {
				throw new Error(
					`incorrect userId response ${this.baseDomain} #${req.status}: ${JSON.stringify(data)}`
				)
			}

			return data.data.id
		}
	}

	async #getUserData() {
		const userId = await this.#getUserId().catch((e) => {
			console.log(e)
			return null
		})

		if (!userId) {
			if (!this.isEu) {
				this.baseDomain = 'donatepay.eu'
				return this.connect()
			}

			throw new Error(`cannot get userId ${this.baseDomain} for ${this.twitchUserId}`)
		}

		while (true) {
			const { isAllowed } = await requestLimiter.consume('donatepay')
			if (!isAllowed) {
				await sleep(1000)
				continue
			}

			const req = await fetch(`https://${this.baseDomain}/api/v2/socket/token`, {
				method: 'post',
				body: JSON.stringify({
					access_token: this.apiKey,
				}),
				headers: {
					'Content-Type': 'application/json',
				},
			})
			if (req.status === 429) {
				await sleep(1000)
				continue
			}

			if (!req.ok) {
				console.error(`DonatePay: cannot get token for ${this.twitchUserId}`, await req.text())
				throw new Error(`cannot get token for ${this.twitchUserId}`)
			}

			const data = await req.json()

			return {
				token: data.token,
				id: userId,
			}
		}
	}
}

interface DonatePayEvent {
	data: {
		notification: {
			type: 'donation'
			vars: {
				name: string
				comment: string
				sum: number
				currency: 'string'
			}
		}
	}
}
