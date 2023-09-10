import * as IO from 'socket.io-client';

import { removeIntegration } from '../index.js';
import { db } from '../libs/db.js';
import { Integration } from '../types.js';
import { onDonation } from '../utils/onDonation.js';

type Socket = typeof IO.Socket;

export class StreamLabs {
	#conn: Socket | null;

	constructor(token: string, private readonly twitchUserId: string) {
		this.#conn = IO.connect(`https://sockets.streamlabs.com?token=${token}`, {
			transports: ['websocket'],
		});

		this.#conn.on('event', (eventData: Event) => {
			if (eventData.type === 'donation') {
				eventData.message.forEach((m) => {
					onDonation({
						twitchUserId: this.twitchUserId,
						amount: m.amount,
						currency: m.currency,
						message: m.message,
						userName: m.from,
					});
				});
			}
		});
	}

	async destroy() {
		this.#conn!.close();
		this.#conn = null;
	}
}

export type Event = {
	type: 'donation';
	message: Message[];
	for: string;
	event_id: string;
};

export type Message = {
	name: string;
	isTest: boolean;
	formatted_amount: string;
	amount: number;
	message: string | null;
	currency: string;
	to: { name: string };
	from: string;
	from_user_id: number;
	_id: string;
	priority: number;
};

export async function addStreamlabsIntegration(integration: Integration) {
	if (
		!integration.accessToken ||
		!integration.refreshToken ||
		!integration.integration ||
		!integration.integration.clientId ||
		!integration.integration.clientSecret ||
		!integration.integration.redirectUrl
	) {
		return;
	}

	await removeIntegration(integration);

	const refresh = await fetch('https://www.twitchalerts.com/api/v1.0/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded',
		},
		body: new URLSearchParams({
			grant_type: 'refresh_token',
			refresh_token: integration.refreshToken,
			redirect_url: integration.integration.redirectUrl,
			client_id: integration.integration.clientId,
			client_secret: integration.integration.clientSecret,
		}).toString(),
	});

	if (!refresh.ok) {
		console.error(await refresh.text());
		return;
	}

	const refreshResponse = await refresh.json();

	await db('channels_integrations').where('id', integration.id).update({
		accessToken: refreshResponse.access_token,
		refreshToken: refreshResponse.refresh_token,
	});

	const socketRequest = await fetch(
		`https://streamlabs.com/api/v1.0/socket/token?access_token=${refreshResponse.access_token}`,
	);

	if (!socketRequest.ok) {
		console.error(await socketRequest.text());
		return;
	}

	const { socket_token } = await socketRequest.json();

	const instance = new StreamLabs(socket_token, integration.channelId);

	return instance;
}
