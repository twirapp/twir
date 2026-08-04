export interface ProviderTokens {
	readonly accessToken: string
	readonly refreshToken: string
}

export interface ProviderTokenStore {
	getTokens(channelID: string): Promise<ProviderTokens | null>
	updateTokens(channelID: string, tokens: ProviderTokens): Promise<void>
}

export interface TokenQuery {
	<Row extends Record<string, unknown>>(
		strings: TemplateStringsArray,
		...values: readonly unknown[]
	): Promise<Row[]>
}

interface TokenRow extends Record<string, unknown> {
	readonly accessToken: string | null
	readonly refreshToken: string | null
}

interface UpdatedRow extends Record<string, unknown> {
	readonly id: string
}

export class ProviderTokensNotFoundError extends Error {
	constructor() {
		super('Enabled provider integration was not found')
		this.name = 'ProviderTokensNotFoundError'
	}
}

export class InvalidProviderTokensError extends Error {
	constructor() {
		super('Provider access and refresh tokens must not be blank')
		this.name = 'InvalidProviderTokensError'
	}
}

function mapTokens(row: TokenRow | undefined): ProviderTokens | null {
	if (!row?.accessToken || !row.refreshToken) {
		return null
	}
	return {
		accessToken: row.accessToken,
		refreshToken: row.refreshToken,
	}
}

function assertUpdated(rows: readonly UpdatedRow[]): void {
	if (rows.length === 0) {
		throw new ProviderTokensNotFoundError()
	}
}

function validateTokens(tokens: ProviderTokens): void {
	if (!tokens.accessToken.trim() || !tokens.refreshToken.trim()) {
		throw new InvalidProviderTokensError()
	}
}

export function createProviderTokenStores(query: TokenQuery): {
	readonly streamElements: ProviderTokenStore
	readonly streamLabs: ProviderTokenStore
} {
	const streamElements: ProviderTokenStore = {
		async getTokens(channelID) {
			const rows = await query<TokenRow>`
				SELECT ci."accessToken" AS "accessToken", ci."refreshToken" AS "refreshToken"
				FROM channels_integrations AS ci
				INNER JOIN integrations AS i ON i.id = ci."integrationId"
				WHERE ci."channelId" = ${channelID}
					AND ci.enabled = TRUE
					AND i.service = 'STREAMELEMENTS'
				LIMIT 1
			`
			return mapTokens(rows[0])
		},
		async updateTokens(channelID, tokens) {
			validateTokens(tokens)
			const rows = await query<UpdatedRow>`
				UPDATE channels_integrations AS ci
				SET "accessToken" = ${tokens.accessToken}, "refreshToken" = ${tokens.refreshToken}
				FROM integrations AS i
				WHERE ci."integrationId" = i.id
					AND ci."channelId" = ${channelID}
					AND ci.enabled = TRUE
					AND i.service = 'STREAMELEMENTS'
				RETURNING ci.id
			`
			assertUpdated(rows)
		},
	}

	const streamLabs: ProviderTokenStore = {
		async getTokens(channelID) {
			const rows = await query<TokenRow>`
				SELECT access_token AS "accessToken", refresh_token AS "refreshToken"
				FROM channels_integrations_streamlabs
				WHERE channel_id = ${channelID} AND enabled = TRUE
				LIMIT 1
			`
			return mapTokens(rows[0])
		},
		async updateTokens(channelID, tokens) {
			validateTokens(tokens)
			const rows = await query<UpdatedRow>`
				UPDATE channels_integrations_streamlabs
				SET access_token = ${tokens.accessToken}, refresh_token = ${tokens.refreshToken}
				WHERE channel_id = ${channelID} AND enabled = TRUE
				RETURNING id
			`
			assertUpdated(rows)
		},
	}

	return { streamElements, streamLabs }
}
