import Centrifuge from 'centrifuge';
import WebSocket from 'ws';

import { onDonation } from '../utils/onDonation.js';

export class DonationAlerts {
	/**
	 * @type {Centrifuge | null}
	 */
	#socket;
	/**
	 *
	 * @type {Centrifuge.Subscription | null}
	 */
	#channel;

	#accessToken;
	#donationAlertsUserId;
	#socketConnectionToken;
	#twitchUserId;

	/**
	 *
	 * @param {string} accessToken
	 * @param {string} donationAlertsUserId
	 * @param {string} socketConnectionToken
	 * @param {string} twitchUserId
	 */
	constructor(
		accessToken,
		donationAlertsUserId,
		socketConnectionToken,
		twitchUserId,
	) {
		this.#accessToken = accessToken;
		this.#donationAlertsUserId = donationAlertsUserId;
		this.#socketConnectionToken = socketConnectionToken;
		this.#twitchUserId = twitchUserId;
	}

	async init() {
		this.#socket = new Centrifuge('wss://centrifugo.donationalerts.com/connection/websocket', {
			websocket: WebSocket,
			onPrivateSubscribe: async (ctx, cb) => {
				const request = await fetch('https://www.donationalerts.com/api/v1/centrifuge/subscribe', {
					method: 'POST',
					body: JSON.stringify(ctx.data),
					headers: { Authorization: `Bearer ${this.#accessToken}` },
				});

				const response = await request.json();
				if (!request.ok) {
					console.error(response);
					cb({ status: request.status, data: {} });
				}

				cb({ status: 200, data: { channels: response.channels } });
			},
		});

		this.#socket.setToken(this.#socketConnectionToken);
		this.#socket.connect();
		
		this.#socket.on('connect', () => {
			console.info(`Connected to donationAlerts #${this.#donationAlertsUserId}`);
		});

		this.#channel = this.#socket.subscribe(`$alerts:donation_${this.#donationAlertsUserId}`);

		this.#channel.on('publish', async ({ data }) => this.#donateCallback(data));

		return this;
	}

	/**
	 * @param {DonationAlertsMessage} data
	 */
	async #donateCallback(data) {
		await onDonation({
			twitchUserId: this.#twitchUserId,
			amount: data.amount,
			currency: data.currency,
			message: data.message,
			userName: data.username,
		});
	}

	async destroy() {
		this.#channel?.removeAllListeners()?.unsubscribe();
		this.#socket?.removeAllListeners()?.disconnect();

		this.#socket = null;
		this.#channel = null;
	}
}
