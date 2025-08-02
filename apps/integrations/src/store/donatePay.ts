import { DonatePay } from '../services/donatepay.ts'

import type { DonatePayIntegration } from '../libs/db.ts'

export const donatePayStore = new Map<string, DonatePay>()

export async function addIntegration(integration: DonatePayIntegration) {
	if (!integration.api_key || !integration.enabled) {
		return
	}

	if (donatePayStore.get(integration.channel_id)) {
		await removeIntegration(integration.channel_id)
	}

	const instance = new DonatePay(integration.channel_id, integration.api_key)
	await instance.connect()

	donatePayStore.set(integration.channel_id, instance)

	return instance
}

export async function removeIntegration(channelId: string) {
	const existed = donatePayStore.get(channelId)
	if (!existed) return

	await existed.disconnect()
	donatePayStore.delete(channelId)
}
