import { db } from '../libs/db.js';
import { Nightbot } from '../services/nightbot.js';

/**
 *
 * @type {Map<string, Nightbot>}
 */
export const nightbotStore = new Map();

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

	removeIntegration(integration.channelId);

	const refresh = await fetch('https://api.nightbot.tv/oauth2/token', {
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

	const instance = new Nightbot(refreshResponse.access_token);

	return instance;
};

/**
 *
 * @param {string} channelId
 */
export const removeIntegration = async (channelId) => {
	const existed = nightbotStore.get(channelId);
	if (!existed) return;
	await existed.destroy();
	nightbotStore.delete(channelId);
};
