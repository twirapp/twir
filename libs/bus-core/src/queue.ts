import { JSONCodec } from 'nats'

import type { NatsConnection, Subscription } from 'nats'

const codec = JSONCodec()

export class Queue<Req, Res> {
	#nc: NatsConnection
	#subscription: Subscription
	private readonly subject: string

	constructor(natsConn: NatsConnection, subject: string) {
		this.#nc = natsConn
		this.subject = subject
	}

	async request(data: Req): Promise<QueueResponse<Res>> {
		const req = await this.#nc.request(this.subject, codec.encode(data))
		return {
			data: codec.decode(req.data) as Res,
		}
	}

	async publish(data: Req): Promise<void> {
		this.#nc.publish(this.subject, codec.encode(data))
	}

	async subscribeGroup(queue: string, callback: QueueSubscribeCallback<Req, Res>) {
		this.#subscription = this.#nc.subscribe(this.subject, { queue })
		for await (const msg of this.#subscription) {
			try {
				const response = await callback(codec.decode(msg.data) as Req)
				if (msg.reply) {
					this.#nc.publish(msg.reply, codec.encode(response))
				}
			} catch (err) {
				const callbackError = err instanceof Error ? err : new Error(String(err))
				console.error(`[bus] subscribeGroup(${queue}): callback error on ${this.subject}:`, callbackError)
				if (msg.reply) {
					try {
						this.#nc.publish(msg.reply, codec.encode({ error: callbackError.message }))
					} catch (publishErr) {
						const replyError = publishErr instanceof Error ? publishErr : new Error(String(publishErr))
						console.error(`[bus] subscribeGroup(${queue}): failed to send error reply:`, replyError)
					}
				}
			}
		}
		console.warn(`[bus] subscribeGroup(${queue}): subscription loop exited unexpectedly on ${this.subject}`)
	}

	async subscribe(callback: QueueSubscribeCallback<Req, Res>) {
		this.#subscription = this.#nc.subscribe(this.subject)
		for await (const msg of this.#subscription) {
			callback(codec.decode(msg.data) as Req)
		}
	}

	async unsubscribe() {
		this.#subscription.unsubscribe()
	}
}

export interface QueueResponse<T> {
	data: T
}

export type QueueSubscribeCallback<Req, Res> = (data: Req) => Res | Promise<Res>
