import { DonateX } from '../services/donatex.ts'

import type { DonateXIntegration } from '../libs/db.ts'

export const donateXStore = new Map<string, DonateX>()

export async function addIntegration(integration: DonateXIntegration) {
	if (!integration.access_token || !integration.refresh_token || !integration.enabled) {
		return
	}

	if (donateXStore.get(integration.channel_id)) {
		await removeIntegration(integration.channel_id)
	}

	const instance = new DonateX(
		integration.channel_id,
		integration.access_token,
		integration.refresh_token
	)

	try {
		await instance.init()
	} catch (e) {
		console.error(e)
	}

	donateXStore.set(integration.channel_id, instance)

	return instance
}

export async function removeIntegration(channelId: string) {
	const existed = donateXStore.get(channelId)
	if (!existed) return

	await existed.destroy()
	donateXStore.delete(channelId)
}
