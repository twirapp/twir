import { sleep } from 'bun'

import {
	DonationAlertsIntegration,
	updateDonationAlertsIntegration,
	updateIntegration,
} from '../libs/db.ts'
import { DonationAlerts, globalRequestLimiter, rateLimiterKey } from '../services/donationAlerts.ts'
import { config } from '@twir/config'

export const donationAlertsStore = new Map<string, DonationAlerts>()

export async function addIntegration(integration: DonationAlertsIntegration) {
	if (!integration.access_token || !integration.refresh_token || !integration.enabled) {
		return
	}

	if (donationAlertsStore.get(integration.channel_id)) {
		await removeIntegration(integration.channel_id)
	}

	let accessToken
	let refreshToken

	while (true) {
		const { isAllowed } = await globalRequestLimiter.consume(rateLimiterKey)
		if (!isAllowed) {
			await sleep(1000)
			continue
		}

		try {
			const refresh = await fetch('https://www.donationalerts.com/oauth/token', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/x-www-form-urlencoded',
				},
				body: new URLSearchParams({
					grant_type: 'refresh_token',
					refresh_token: integration.refresh_token,
					client_id: config.DONATIONALERTS_CLIENT_ID!,
					client_secret: config.DONATIONALERTS_CLIENT_SECRET!,
				}).toString(),
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
		} catch (e) {
			console.log(e)
			await sleep(1000)
			continue
		}
	}

	if (!accessToken || !refreshToken) {
		return
	}

	await updateDonationAlertsIntegration({
		channel_id: integration.channel_id,
		access_token: accessToken,
		refresh_token: refreshToken,
	})

	let profileData

	while (true) {
		const { isAllowed } = await globalRequestLimiter.consume(rateLimiterKey)
		if (!isAllowed) {
			await sleep(1000)
			continue
		}

		const request = await fetch('https://www.donationalerts.com/api/v1/user/oauth', {
			headers: {
				Authorization: `Bearer ${accessToken}`,
			},
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
		integration.channel_id
	)
	await instance.init()

	donationAlertsStore.set(integration.channel_id, instance)

	return instance
}

export async function removeIntegration(channelId: string) {
	const existed = donationAlertsStore.get(channelId)
	if (!existed) return

	await existed.destroy()
	donationAlertsStore.delete(channelId)
}
