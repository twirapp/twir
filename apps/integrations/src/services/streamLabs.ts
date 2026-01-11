import io from 'socket.io-client'

import { onDonation } from '../utils/onDonation.js'

export class StreamLabs {
	#conn: typeof io.Socket | null

	constructor(
		private readonly token: string,
		private readonly twitchUserId: string
	) {
		this.#conn = io.connect(`https://sockets.streamlabs.com?token=${token}`, {
			transports: ['websocket'],
		})

		this.#conn!.on('connect', () => {
			console.log(`StreamLabs connected for user ${this.twitchUserId}`)
		})
		this.#conn!.on('event', (eventData: StreamLabsEvent) => this.#eventCallback(eventData))
	}

	#eventCallback(event: StreamLabsEvent) {
		if (event.type === 'donation') {
			event.message.forEach((m) => {
				onDonation({
					twitchUserId: this.twitchUserId,
					amount: m.amount,
					currency: m.currency,
					message: m.message,
					userName: m.from,
				})
			})
		}
	}

	async destroy() {
		this.#conn?.close()
		this.#conn = null
	}
}

export interface StreamLabsEvent {
	type: 'donation'
	message: StreamLabsMessage[]
	for: string
	event_id: string
}

export interface StreamLabsMessage {
	name: string
	isTest: boolean
	formatted_amount: string
	amount: number
	message: string | null
	currency: string
	to: { name: string }
	from: string
	from_user_id: number
	_id: string
	priority: number
}
