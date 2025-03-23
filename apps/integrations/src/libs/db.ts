import process from 'node:process'

import { config } from '@twir/config'
import { SQL } from 'bun'

import type { Donate } from '../utils/onDonation.ts'

const sql = new SQL(config.DATABASE_URL)

try {
	await sql`SELECT 1`
	console.log('Connected to database')
} catch (e) {
	console.error(e)
	process.exit(1)
}

export const Service = Object.freeze({
	DONATIONALERTS: 'DONATIONALERTS',
	STREAMLABS: 'STREAMLABS',
	DONATEPAY: 'DONATEPAY',
})

export async function getIntegrations(integrationId: string): Promise<Integration | null>
export async function getIntegrations(): Promise<Integration[]>
export async function getIntegrations(integrationId?: string): Promise<Integration | Integration[] | null> {
	const result = await sql`
	SELECT channel_integration.id,
	       channel_integration.enabled,
	       channel_integration."accessToken",
	       channel_integration."refreshToken",
	       channel_integration."clientId",
	       channel_integration."clientSecret",
	       channel_integration."apiKey",
	       channel_integration.data,
	       channel_integration."channelId",
	       channel_integration."integrationId",
				 (json_agg(integration.*) ->> 0)::json as integration,
				 (json_agg(channel.*) ->> 0)::json as channel
	FROM channels_integrations channel_integration
	LEFT JOIN integrations integration ON integration.id = channel_integration."integrationId"
	LEFT JOIN channels channel ON channel.id = channel_integration."channelId"
	WHERE channel_integration.enabled = true AND integration.service IN ('DONATIONALERTS', 'STREAMLABS', 'DONATEPAY')
		${
	integrationId
		? sql`
	AND channel_integration.id = ${integrationId} GROUP BY channel_integration.id, integration.id, channel.id
`
		: sql`GROUP BY channel_integration.id, integration.id, channel.id`
}

`

	if (integrationId) {
		return result[0] ?? null
	}

	return result || []
}

export async function updateIntegration(id: string, data: {
	enabled?: boolean
	accessToken?: string
	refreshToken?: string
}) {
	if (Object.keys(data).length === 0) {
		return
	}

	await sql.begin(async tx => {
		if (data.enabled !== undefined) {
			await tx`UPDATE channels_integrations SET "enabled" = ${data.enabled} WHERE id = ${id}`
		}

		if (data.accessToken !== undefined) {
			await tx`UPDATE channels_integrations SET "accessToken" = ${data.accessToken} WHERE id = ${id}`
		}

		if (data.refreshToken !== undefined) {
			await tx`UPDATE channels_integrations SET "refreshToken" = ${data.refreshToken} WHERE id = ${id}`
		}
	})
}

export async function insertDonation(data: Donate) {
	console.log(data)
	const preparedData = {
		channel_id: data.twitchUserId,
		type: 'DONATION',
		data: {
			donationAmount: data.amount.toString(),
			donationCurrency: data.currency,
			donationMessage: data.message,
			donationUsername: data.userName ?? 'Anonymous',
		},
	}

	await sql`
	INSERT INTO channels_events_list
	${sql(preparedData)}
`
}

export interface Integration {
	id: string
	enabled: true
	accessToken: string | null
	refreshToken: string | null
	clientId: string | null
	clientSecret: string | null
	apiKey: string | null
	data: Record<string, any> | null
	channelId: string
	integrationId: string
	integration: {
		id: string
		service: keyof typeof Service
		accessToken: string | null
		refreshToken: string | null
		clientId: string | null
		clientSecret: string | null
		apiKey: string | null
		redirectUrl: string | null
	}
	channel: {
		id: string
		isEnabled: boolean
		isTwitchBanned: boolean
		isBanned: boolean
		botId: string
		isBotMod: boolean
	}
}
