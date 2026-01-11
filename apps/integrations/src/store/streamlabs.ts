import { config } from '@twir/config'

import { StreamlabsIntegration, updateStreamlabsIntegration } from '../libs/db.ts'
import { StreamLabs } from '../services/streamLabs.ts'

export const streamLabsStore = new Map<string, StreamLabs>()

export async function addIntegration(integration: StreamlabsIntegration) {
	if (!integration.access_token || !integration.refresh_token || !integration.enabled) {
		return
	}

	await removeIntegration(integration.channel_id)

	const refresh = await fetch('https://streamlabs.com/api/v2.0/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Accept: 'application/json',
		},
		body: JSON.stringify({
			grant_type: 'refresh_token',
			refresh_token: integration.refresh_token,
			redirect_url: `${config.SITE_BASE_URL}/dashboard/integrations/streamlabs`,
			client_id: config.STREAMLABS_CLIENT_ID!,
			client_secret: config.STREAMLABS_CLIENT_SECRET!,
		}),
	})

	if (!refresh.ok) {
		console.error(`token request error:`, await refresh.text())
		return
	}

	const refreshResponse = await refresh.json()

	await updateStreamlabsIntegration({
		channel_id: integration.channel_id,
		access_token: refreshResponse.access_token,
		refresh_token: refreshResponse.refresh_token,
	})

	const socketRequest = await fetch(`https://streamlabs.com/api/v2.0/socket/token`, {
		headers: {
			Authorization: `Bearer ${refreshResponse.access_token}`,
			Accept: 'application/json',
		},
	})

	if (!socketRequest.ok) {
		console.error(await socketRequest.text())
		await updateStreamlabsIntegration({
			channel_id: integration.channel_id,
			enabled: 'false',
		})
		return
	}

	const { socket_token } = await socketRequest.json()

	return new StreamLabs(socket_token, integration.channel_id)
}

export async function removeIntegration(channelId: string) {
	const existed = streamLabsStore.get(channelId)
	if (!existed) return
	await existed.destroy()
	streamLabsStore.delete(channelId)
}
