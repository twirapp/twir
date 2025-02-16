import { updateIntegration } from '../libs/db.ts'
import { StreamLabs } from '../services/streamLabs.ts'

import type { Integration } from '../libs/db.ts'

export const streamLabsStore = new Map<string, StreamLabs>()

export async function addIntegration(integration: Integration) {
	if (
		!integration.accessToken
		|| !integration.refreshToken
		|| !integration.integration
		|| !integration.integration.clientId
		|| !integration.integration.clientSecret
		|| !integration.integration.redirectUrl
		|| !integration.enabled
	) {
		return
	}

	await removeIntegration(integration.channelId)

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
	})

	if (!refresh.ok) {
		console.error(await refresh.text())
		return
	}

	const refreshResponse = await refresh.json()

	await updateIntegration(integration.id, {
		accessToken: refreshResponse.access_token,
		refreshToken: refreshResponse.refresh_token,
	})

	const socketRequest = await fetch(
		`https://streamlabs.com/api/v1.0/socket/token?access_token=${refreshResponse.access_token}`,
	)

	if (!socketRequest.ok) {
		console.error(await socketRequest.text())
		await updateIntegration(integration.id, { enabled: false })
		return
	}

	const { socket_token } = await socketRequest.json()

	return new StreamLabs(socket_token, integration.channelId)
}

export async function removeIntegration(channelId: string) {
	const existed = streamLabsStore.get(channelId)
	if (!existed) return
	await existed.destroy()
	streamLabsStore.delete(channelId)
}
