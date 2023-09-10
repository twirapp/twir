import Centrifuge from 'centrifuge';
import WebSocket from 'ws';

import { donationAlertsStore, removeIntegration } from '../index.js';
import { db } from '../libs/db.js';
import { Integration } from '../types.js';
import { onDonation } from '../utils/onDonation.js';

export class DonationAlerts {
	#socket: Centrifuge | null;
	#channel: Centrifuge.Subscription | null;

	constructor(
		private readonly accessToken: string,
		private readonly donationAlertsUserId: string,
		private readonly socketConnectionToken: string,
		private readonly twitchUserId: string,
	) {
	}

	async init() {
		this.#socket = new Centrifuge('wss://centrifugo.donationalerts.com/connection/websocket', {
			websocket: WebSocket,
			onPrivateSubscribe: async (ctx, cb) => {
				const request = await fetch('https://www.donationalerts.com/api/v1/centrifuge/subscribe', {
					method: 'POST',
					body: JSON.stringify(ctx.data),
					headers: { Authorization: `Bearer ${this.accessToken}` },
				});

				const response = await request.json();
				if (!request.ok) {
					console.error(response);
					cb({ status: request.status, data: {} as any });
				}

				cb({ status: 200, data: { channels: response.channels } });
			},
		});

		this.#socket.setToken(this.socketConnectionToken);
		this.#socket.connect();

		this.#channel = this.#socket.subscribe(`$alerts:donation_${this.donationAlertsUserId}`);

		this.#channel.on('publish', async ({ data }: { data: Message }) => {
			await onDonation({
				twitchUserId: this.twitchUserId,
				amount: data.amount,
				currency: data.currency,
				message: data.message,
				userName: data.username,
			});
		});

		return this;
	}

	async destroy() {
		this.#channel?.removeAllListeners()?.unsubscribe();
		this.#socket?.removeAllListeners()?.disconnect();

		this.#socket = null;
		this.#channel = null;
	}
}

export type Message = {
	id: number;
	name: string;
	username?: string | null;
	message: string | null;
	message_type: 'text' | 'audio';
	payin_system: null | any;
	amount: number;
	currency: string;
	amount_in_user_currency: number;
	recipient_name: string;
	recipient: {
		user_id: number;
		code: string;
		name: string;
		avatar: string;
	};
	created_at: string;
	shown_at: null | any;
	reason: string;
};

export async function addDonationAlertsIntegration(integration: Integration) {
	if (
		!integration.accessToken ||
		!integration.refreshToken ||
		!integration.integration ||
		!integration.integration.clientId ||
		!integration.integration.clientSecret
	) {
		return;
	}

	if (donationAlertsStore.get(integration.channelId)) {
		await removeIntegration(integration);
	}

	const refresh = await fetch('https://www.donationalerts.com/oauth/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded',
		},
		body: new URLSearchParams({
			grant_type: 'refresh_token',
			refresh_token: integration.refreshToken,
			client_id: integration.integration.clientId,
			client_secret: integration.integration.clientSecret,
		}).toString(),
	});

	if (!refresh.ok) {
		console.error('cannot refresh DA tokens:', await refresh.text());
		return;
	}

	const refreshResponse = await refresh.json();

	await db('channels_integrations').where('id', integration.id).update({
		accessToken: refreshResponse.access_token,
		refreshToken: refreshResponse.refresh_token,
	});

	const request = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
		headers: {
			Authorization: `Bearer ${refreshResponse.access_token}`,
		},
	});

	if (!request.ok) {
		console.log(await request.text());
		return;
	}

	const { data } = await request.json();
	const { id, socket_connection_token } = data;
	const instance = new DonationAlerts(
		refreshResponse.access_token,
		id,
		socket_connection_token,
		integration.channelId,
	);
	await instance.init();

	return instance;
}
