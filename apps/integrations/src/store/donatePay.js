import { DonatePay } from '../services/donatepay.js';

/**
 *
 * @type {Map<string, DonatePay>}
 */
export const donatePayStore = new Map();

/**
 *
 * @param {Integration} integration
 */
export async function addIntegration(integration) {
	if (
		!integration.integration ||
		!integration.apiKey
	) {
		return;
	}

	if (donatePayStore.get(integration.channelId)) {
		await removeIntegration(integration);
	}

	const instance = new DonatePay(integration.channelId, integration.apiKey);
	await instance.connect();

	donatePayStore.set(integration.channelId, instance);

	return instance;
}

/**
 *
 * @param channelId
 */
export const removeIntegration = async (channelId) => {
	const existed = donatePayStore.get(channelId);
	if (!existed) return;

	await existed.destroy();
	donatePayStore.delete(channelId);
};
