import { db } from '../libs/db.js';
import { StreamLabs } from '../services/streamLabs.js';

/**
 *
 * @type {Map<string, StreamLabs>}
 */
export const streamlabsStore = new Map();

/**
 *
 * @param {Integration} integration
 */
export const addIntegration = async (integration) => {
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

	await removeIntegration(integration.channelId);

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
};

/**
 *
 * @param {string} channelId
 */
export const removeIntegration = async (channelId) => {
	const existed = streamlabsStore.get(channelId);
	if (!existed) return;
	await existed.destroy();
	streamlabsStore.delete(channelId);
};
