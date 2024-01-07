import * as IO from 'socket.io-client';

import { onDonation } from '../utils/onDonation.js';

export class StreamLabs {
	/**
	 * @private
	 * @type {IO.Socket | null}
	 */
	#conn;

	/**
	 *
	 * @param {string} token
	 * @param {string} twitchUserId
	 */

	#twitchUserId;
	constructor(token, twitchUserId) {
		this.#twitchUserId = twitchUserId;

		this.#conn = IO.connect(`https://sockets.streamlabs.com?token=${token}`, {
			transports: ['websocket'],
		});

		this.#conn.on('event', (eventData) => this.#eventCallback(eventData));
	}

	/**
	 *
	 * @param {StreamLabsEvent} event
	 */
	#eventCallback(event) {
		if (event.type === 'donation') {
			event.message.forEach((m) => {
				onDonation({
					twitchUserId: this.#twitchUserId,
					amount: m.amount,
					currency: m.currency,
					message: m.message,
					userName: m.from,
				});
			});
		}
	}

	async destroy() {
		this.#conn.close();
		this.#conn = null;
	}
}
