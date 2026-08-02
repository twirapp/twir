import { withOAuthRefreshLock } from './oauth-lock.ts'
import type { ProviderTokenStore, ProviderTokens } from './provider-token-store.ts'

const API_BASE_URL = 'https://streamlabs.com'
const REQUEST_TIMEOUT_MS = 15_000
const MAX_RESPONSE_BYTES = 1024 * 1024

export type StreamLabsFetch = (
	input: Parameters<typeof globalThis.fetch>[0],
	init?: Parameters<typeof globalThis.fetch>[1]
) => Promise<Response>

export type StreamLabsRefreshLock = <T>(
	provider: 'streamlabs',
	channelID: string,
	callback: (signal: AbortSignal) => Promise<T>
) => Promise<T>

export interface StreamLabsClientOptions {
	readonly channelID: string
	readonly tokens: ProviderTokens
	readonly tokenStore: ProviderTokenStore
	readonly clientID: string
	readonly clientSecret: string
	readonly redirectURI: string
	readonly fetch?: StreamLabsFetch
	readonly lock?: StreamLabsRefreshLock
	readonly requestTimeoutMs?: number
}

interface SocketTokenResponse {
	readonly socket_token: string
}

export interface StreamLabsSocketToken {
	readonly socketToken: string
}

interface RefreshResponse {
	readonly access_token: string
	readonly refresh_token?: string
}

class UnauthorizedError extends Error {}

export class StreamLabsClientError extends Error {
	constructor() {
		super('Streamlabs provider request failed')
		this.name = 'StreamLabsClientError'
	}
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return value !== null && typeof value === 'object'
}

function isSocketTokenResponse(value: unknown): value is SocketTokenResponse {
	return isRecord(value)
		&& typeof value.socket_token === 'string'
		&& value.socket_token.trim().length > 0
}

function isRefreshResponse(value: unknown): value is RefreshResponse {
	return isRecord(value)
		&& typeof value.access_token === 'string'
		&& value.access_token.trim().length > 0
		&& (value.refresh_token === undefined
			|| (typeof value.refresh_token === 'string' && value.refresh_token.trim().length > 0))
}

async function readBoundedJson(response: Response): Promise<unknown> {
	const contentLength = Number(response.headers.get('content-length'))
	if (Number.isFinite(contentLength) && contentLength > MAX_RESPONSE_BYTES) {
		throw new Error('Streamlabs response is too large')
	}
	if (!response.body) return null

	const reader = response.body.getReader()
	const chunks: Uint8Array[] = []
	let bytes = 0
	while (true) {
		const result = await reader.read()
		if (result.done) break
		bytes += result.value.byteLength
		if (bytes > MAX_RESPONSE_BYTES) {
			await reader.cancel()
			throw new Error('Streamlabs response is too large')
		}
		chunks.push(result.value)
	}

	const body = new Uint8Array(bytes)
	let offset = 0
	for (const chunk of chunks) {
		body.set(chunk, offset)
		offset += chunk.byteLength
	}
	try {
		return JSON.parse(new TextDecoder().decode(body)) as unknown
	} catch {
		throw new Error('Streamlabs response was malformed')
	}
}

export class StreamLabsClient {
	readonly #channelID: string
	readonly #tokenStore: ProviderTokenStore
	readonly #clientID: string
	readonly #clientSecret: string
	readonly #redirectURI: string
	readonly #fetch: StreamLabsFetch
	readonly #lock: StreamLabsRefreshLock
	readonly #requestTimeoutMs: number
	#tokens: ProviderTokens

	constructor(options: StreamLabsClientOptions) {
		this.#channelID = options.channelID
		this.#tokens = options.tokens
		this.#tokenStore = options.tokenStore
		this.#clientID = options.clientID
		this.#clientSecret = options.clientSecret
		this.#redirectURI = options.redirectURI
		this.#fetch = options.fetch ?? ((input, init) => globalThis.fetch(input, init))
		this.#lock = options.lock ?? withOAuthRefreshLock
		this.#requestTimeoutMs = options.requestTimeoutMs ?? REQUEST_TIMEOUT_MS
	}

	get tokens(): ProviderTokens {
		return { ...this.#tokens }
	}

	async getSocketToken(): Promise<StreamLabsSocketToken> {
		try {
			const attemptedAccessToken = this.#tokens.accessToken
			try {
				return await this.#requestSocketToken(attemptedAccessToken)
			} catch (error) {
				if (!(error instanceof UnauthorizedError)) throw error
			}

			await this.#refreshAfterUnauthorized(attemptedAccessToken)
			return await this.#requestSocketToken(this.#tokens.accessToken)
		} catch {
			throw new StreamLabsClientError()
		}
	}

	async #requestSocketToken(accessToken: string): Promise<StreamLabsSocketToken> {
		const response = await this.#fetch(`${API_BASE_URL}/api/v2.0/socket/token`, {
			headers: {
				Accept: 'application/json',
				Authorization: `Bearer ${accessToken}`,
			},
			signal: AbortSignal.timeout(this.#requestTimeoutMs),
		})
		if (response.status === 401) throw new UnauthorizedError()
		if (!response.ok) throw new Error(`Streamlabs socket token returned status ${response.status}`)
		const decoded = await readBoundedJson(response)
		if (!isSocketTokenResponse(decoded)) throw new Error('Streamlabs socket token response was invalid')
		return { socketToken: decoded.socket_token }
	}

	async #refreshAfterUnauthorized(attemptedAccessToken: string): Promise<void> {
		await this.#lock('streamlabs', this.#channelID, async (lockSignal) => {
			const persisted = await this.#tokenStore.getTokens(this.#channelID)
			if (!persisted) throw new Error('Enabled Streamlabs integration was not found')
			if (persisted.accessToken !== attemptedAccessToken) {
				this.#tokens = persisted
				return
			}

			const response = await this.#fetch(`${API_BASE_URL}/api/v2.0/token`, {
				method: 'POST',
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/x-www-form-urlencoded',
				},
				body: new URLSearchParams({
					grant_type: 'refresh_token',
					client_id: this.#clientID,
					client_secret: this.#clientSecret,
					refresh_token: persisted.refreshToken,
					redirect_uri: this.#redirectURI,
				}),
				signal: AbortSignal.any([
					lockSignal,
					AbortSignal.timeout(this.#requestTimeoutMs),
				]),
			})
			if (!response.ok) throw new Error(`Streamlabs refresh returned status ${response.status}`)
			const decoded = await readBoundedJson(response)
			if (!isRefreshResponse(decoded)) throw new Error('Streamlabs refresh response was invalid')
			const nextTokens = {
				accessToken: decoded.access_token,
				refreshToken: decoded.refresh_token ?? persisted.refreshToken,
			}
			await this.#tokenStore.updateTokens(this.#channelID, nextTokens)
			this.#tokens = nextTokens
		})
	}
}
