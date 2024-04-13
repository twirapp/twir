import { setTimeout as sleep } from 'timers/promises';

import { db } from '../libs/db.js';
import { DonationAlerts } from '../services/donationAlerts.js';

/**
 *
 * @type {Map<string, DonationAlerts>}
 */
export const donationAlertsStore = new Map();

/**
 *
 * @param {Integration} integration
 */
export async function addIntegration(integration) {
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
		await removeIntegration(integration.channelId);
	}

	let accessToken;
	let refreshToken;

	while(true) {
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
			if (refresh.status === 429) {
				await sleep(5000);
				continue;
			}
			console.error('cannot refresh DA tokens:', await refresh.text());
			break;
		}

		const refreshResponse = await refresh.json();
		accessToken = refreshResponse.access_token;
		refreshToken = refreshResponse.refresh_token;
		break;
	}

	if (!accessToken || !refreshToken) {
		return;
	}

	await db('channels_integrations').where('id', integration.id).update({
		accessToken: accessToken,
		refreshToken: refreshToken,
	});

	let profileData;

	while(true) {
		const request = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
			headers: {
				Authorization: `Bearer ${accessToken}`,
			},
		});

		if (!request.ok) {
			if (request.status === 429) {
				await sleep(5000);
				continue;
			}
			console.error('cannot get donationAlerts profile', await request.text());
			break;
		}

		const response = await request.json();
		profileData = response.data;
		break;
	}

	const { id, socket_connection_token } = profileData;
	const instance = new DonationAlerts(
		accessToken,
		id,
		socket_connection_token,
		integration.channelId,
	);
	await instance.init();

	donationAlertsStore.set(integration.channelId, instance);

	return instance;
}

/**
 *
 * @param channelId
 */
export const removeIntegration = async (channelId) => {
	const existed = donationAlertsStore.get(channelId);
	if (!existed) return;

	await existed.destroy();
	donationAlertsStore.delete(channelId);
};
