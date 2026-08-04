import type { ProviderTokenStore, ProviderTokens } from './provider-token-store.ts'
import { withOAuthRefreshLock } from './oauth-lock.ts'

const TOKEN_URL = 'https://api.streamelements.com/oauth2/token'
const REQUEST_TIMEOUT_MS = 15_000
const MAX_RESPONSE_BYTES = 1024 * 1024

export type StreamElementsFetch = (
	input: Parameters<typeof globalThis.fetch>[0],
	init?: Parameters<typeof globalThis.fetch>[1]
) => Promise<Response>

export type StreamElementsRefreshLock = <T>(
	provider: 'streamelements',
	channelID: string,
	callback: (signal: AbortSignal) => Promise<T>
) => Promise<T>

interface RefreshResponse {
	readonly access_token: string
	readonly refresh_token?: string
}

export interface StreamElementsClientOptions {
	readonly channelID: string
	readonly tokens: ProviderTokens
	readonly tokenStore: ProviderTokenStore
	readonly clientID: string
	readonly clientSecret: string
	readonly fetch?: StreamElementsFetch
	readonly lock?: StreamElementsRefreshLock
	readonly requestTimeoutMs?: number
}

export class StreamElementsRefreshError extends Error {
	constructor() {
		super('StreamElements token refresh failed')
		this.name = 'StreamElementsRefreshError'
	}
}

function isRefreshResponse(value: unknown): value is RefreshResponse {
	if (!value || typeof value !== 'object') return false
	const response = value as Record<string, unknown>
	return typeof response.access_token === 'string'
		&& response.access_token.trim().length > 0
		&& (response.refresh_token === undefined
			|| (typeof response.refresh_token === 'string' && response.refresh_token.trim().length > 0))
}

async function readBoundedJson(response: Response): Promise<unknown> {
	const contentLength = Number(response.headers.get('content-length'))
	if (Number.isFinite(contentLength) && contentLength > MAX_RESPONSE_BYTES) {
		throw new Error('StreamElements response is too large')
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
			throw new Error('StreamElements response is too large')
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
		throw new Error('StreamElements response was malformed')
	}
}

export class StreamElementsClient {
	readonly #channelID: string
	readonly #tokenStore: ProviderTokenStore
	readonly #clientID: string
	readonly #clientSecret: string
	readonly #fetch: StreamElementsFetch
	readonly #lock: StreamElementsRefreshLock
	readonly #requestTimeoutMs: number
	#tokens: ProviderTokens

	constructor(options: StreamElementsClientOptions) {
		this.#channelID = options.channelID
		this.#tokens = options.tokens
		this.#tokenStore = options.tokenStore
		this.#clientID = options.clientID
		this.#clientSecret = options.clientSecret
		this.#fetch = options.fetch ?? ((input, init) => globalThis.fetch(input, init))
		this.#lock = options.lock ?? withOAuthRefreshLock
		this.#requestTimeoutMs = options.requestTimeoutMs ?? REQUEST_TIMEOUT_MS
	}

	get tokens(): ProviderTokens {
		return { ...this.#tokens }
	}

	async refresh(): Promise<void> {
		try {
			await this.#lock('streamelements', this.#channelID, async (lockSignal) => {
				const persisted = await this.#tokenStore.getTokens(this.#channelID)
				if (!persisted) throw new Error('Enabled StreamElements integration was not found')

				if (persisted.accessToken !== this.#tokens.accessToken) {
					this.#tokens = persisted
					return
				}

				const timeoutSignal = AbortSignal.timeout(this.#requestTimeoutMs)
				const response = await this.#fetch(TOKEN_URL, {
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
					}),
					signal: AbortSignal.any([lockSignal, timeoutSignal]),
				})
				if (!response.ok) throw new Error(`StreamElements refresh returned status ${response.status}`)

				const decoded = await readBoundedJson(response)
				if (!isRefreshResponse(decoded)) throw new Error('StreamElements refresh returned invalid JSON')
				const nextTokens = {
					accessToken: decoded.access_token,
					refreshToken: decoded.refresh_token ?? persisted.refreshToken,
				}
				await this.#tokenStore.updateTokens(this.#channelID, nextTokens)
				this.#tokens = nextTokens
			})
		} catch {
			throw new StreamElementsRefreshError()
		}
	}
}
