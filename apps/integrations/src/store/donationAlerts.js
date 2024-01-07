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
