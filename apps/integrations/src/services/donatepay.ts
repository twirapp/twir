import Centrifuge from 'centrifuge'
import ws from 'ws'
// eslint-disable-next-line ts/ban-ts-comment
// @ts-expect-error
import { XMLHttpRequest } from 'xmlhttprequest'

import { onDonation } from '../utils/onDonation.js'

import type { Subscription } from 'centrifuge'

// eslint-disable-next-line no-restricted-globals
global.XMLHttpRequest = XMLHttpRequest

export class DonatePay {
	#centrifuge: Centrifuge
	#subscription: Subscription
	#timeout: Timer

	constructor(private readonly twitchUserId: string, private readonly apiKey: string) {}

	async connect() {
		if (this.#centrifuge || this.#subscription) {
			await this.disconnect()
		}

		const userData = await this.#getUserData()

		this.#centrifuge = new Centrifuge('wss://centrifugo.donatepay.ru:43002/connection/websocket', {
			subscribeEndpoint: 'https://donatepay.ru/api/v2/socket/token',
			subscribeParams: {
				access_token: this.apiKey,
			},
			disableWithCredentials: true,
			websocket: ws,
			ping: true,
			pingInterval: 5000,
		})

		this.#centrifuge.setToken(userData.token)

		this.#subscription = this.#centrifuge.subscribe(`$public:${userData.id}`, this.#eventCallback)

		const logDisconnect = (args: any[]) => console.info(`DonatePay(${this.twitchUserId}): disconnected`, ...args)

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
		const getUserParams = new URLSearchParams({
			access_token: this.apiKey,
		})
		const req = await fetch(`https://donatepay.ru/api/v1/user?${getUserParams}`)

		if (!req.ok) {
			throw new Error('incorrect response')
		}

		const data = await req.json()

		if (!data.data?.id) {
			throw new Error('incorrect response')
		}

		return data.data.id
	}

	async #getUserData() {
		const userId = await this.#getUserId().catch(() => null)

		if (!userId) {
			console.error(`DonatePay: something went wrong when getting token of ${this.twitchUserId}`)
		}

		const req = await fetch('https://donatepay.ru/api/v2/socket/token', {
			method: 'post',
			body: JSON.stringify({
				access_token: this.apiKey,
			}),
			headers: {
				'Content-Type': 'application/json',
			},
		})
		const data = await req.json()

		return {
			token: data.token,
			id: userId,
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
