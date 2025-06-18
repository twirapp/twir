import { sleep } from 'bun'

import { updateIntegration } from '../libs/db.ts'
import { DonationAlerts, globalRequestLimiter } from '../services/donationAlerts.ts'

import type { Integration } from '../libs/db.ts'

export const donationAlertsStore = new Map<string, DonationAlerts>()

export async function addIntegration(integration: Integration) {
	if (
		!integration.accessToken
		|| !integration.refreshToken
		|| !integration.integration
		|| !integration.integration.clientId
		|| !integration.integration.clientSecret
		|| !integration.enabled
	) {
		return
	}

	if (donationAlertsStore.get(integration.channelId)) {
		await removeIntegration(integration.channelId)
	}

	let accessToken
	let refreshToken

	while (true) {
		const { isAllowed } = await globalRequestLimiter.consume(integration.id)
		if (!isAllowed) {
			await sleep(1000)
			continue
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
			verbose: true,
		})

		if (!refresh.ok) {
			if (refresh.status === 429) {
				await sleep(1000)
				continue
			}
			console.error('cannot refresh DA tokens:', await refresh.text())
			break
		}

		const refreshResponse = await refresh.json()
		accessToken = refreshResponse.access_token
		refreshToken = refreshResponse.refresh_token
		break
	}

	if (!accessToken || !refreshToken) {
		await updateIntegration(integration.id, { enabled: false })
		return
	}

	await updateIntegration(integration.id, { accessToken, refreshToken })

	let profileData

	while (true) {
		const { isAllowed } = await globalRequestLimiter.consume(integration.id)
		if (!isAllowed) {
			await sleep(1000)
			continue
		}

		const request = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
			headers: {
				Authorization: `Bearer ${accessToken}`,
			},
			verbose: true,
		})

		if (!request.ok) {
			if (request.status === 429) {
				await sleep(1000)
				continue
			}
			console.error('cannot get donationAlerts profile', await request.text())
			break
		}

		const response = await request.json()
		profileData = response.data
		break
	}

	if (!profileData) {
		return
	}

	const { id, socket_connection_token } = profileData
	const instance = new DonationAlerts(
		accessToken,
		id,
		socket_connection_token,
		integration.channelId,
	)
	await instance.init()

	donationAlertsStore.set(integration.channelId, instance)

	return instance
}

export async function removeIntegration(channelId: string) {
	const existed = donationAlertsStore.get(channelId)
	if (!existed) return

	await existed.destroy()
	donationAlertsStore.delete(channelId)
}
