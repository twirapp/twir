import { JSONCodec } from 'nats';

const codec = JSONCodec();

export class Queue {
	/**
		* @type {import("nats").NatsConnection} nc
	 */
	#nc;
	/**
	 * @type {import("nats").Subscription} subscription
	 */
	#subscription;

	constructor(natsConn, subject) {
		this.#nc = natsConn;
		this.subject = subject;
	}

	async request(data) {
		const req = await this.#nc.request(this.subject, codec.encode(data));
		return {
			data: codec.decode(req.data),
		};
	}

	async publish(data) {
		await this.#nc.publish(this.subject, codec.encode(data));
	}

	async subscribeGroup(queue, callback) {
		this.#subscription = this.#nc.subscribe(this.subject, { queue });
		for await (const msg of this.#subscription) {
			const response = await callback(codec.decode(msg.data));
			await this.#nc.publish(msg.reply, codec.encode(response));
		}
	}

	async subscribe(callback) {
		this.#subscription = this.#nc.subscribe(this.subject);
		for await (const msg of this.#subscription) {
			callback(codec.decode(msg.data));
		}
	}

	async unsubscribe() {
		await this.#subscription.unsubscribe();
	}
}
