import { DonatePay } from '../services/donatepay.ts'

import type { Integration } from '../libs/db.ts'

export const donatePayStore = new Map<string, DonatePay>()

export async function addIntegration(integration: Integration) {
	if (
		!integration.integration
		|| !integration.apiKey
		|| !integration.enabled
	) {
		return
	}

	if (donatePayStore.get(integration.channelId)) {
		await removeIntegration(integration.channelId)
	}

	const instance = new DonatePay(integration.channelId, integration.apiKey)
	await instance.connect()

	donatePayStore.set(integration.channelId, instance)

	return instance
}

export async function removeIntegration(channelId: string) {
	const existed = donatePayStore.get(channelId)
	if (!existed) return

	await existed.disconnect()
	donatePayStore.delete(channelId)
}
