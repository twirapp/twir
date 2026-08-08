import { sleep } from 'bun'
import { HttpTransportType, HubConnectionBuilder, LogLevel } from '@microsoft/signalr'

import { updateDonateXIntegration } from '../libs/db.ts'
import { onDonation } from '../utils/onDonation'
import { globalRequestLimiter } from './donationAlerts.ts'

import type { HubConnection } from '@microsoft/signalr'
import { config } from '@twir/config'

export const donatexRateLimiterKey = 'donatex'

export class DonateX {
	#connection: HubConnection | null = null
	#refreshPromise: Promise<string> | null = null

	constructor(
		private readonly twitchUserId: string,
		private accessToken: string,
		private refreshToken: string
	) {}

	async init() {
		this.#connection = new HubConnectionBuilder()
			.withUrl('https://donatex.gg/api/public-donations-hub', {
				transport: HttpTransportType.WebSockets,
				// the hub expects the token in the access_token query param,
				// which is exactly what signalr does for websockets
				accessTokenFactory: () => this.#getFreshAccessToken(),
			})
			.withAutomaticReconnect()
			.configureLogging(LogLevel.Information)
			.build()

		this.#connection.on('DonationCreated', (donation: DonateXDonation) =>
			this.#donateCallback(donation)
		)

		this.#connection.onreconnecting((error) => {
			console.warn(`[DONATEX #${this.twitchUserId}] reconnecting: ${error?.message}`)
		})

		this.#connection.onreconnected(() => {
			console.info(`[DONATEX #${this.twitchUserId}] reconnected`)
		})

		await this.#connection.start()
		console.info(`Connected to donatex #${this.twitchUserId}`)

		return this
	}

	#getFreshAccessToken(): Promise<string> {
		// single-flight: refresh tokens rotate, so concurrent refreshes
		// would invalidate each other
		if (!this.#refreshPromise) {
			this.#refreshPromise = this.#refreshAccessToken().finally(() => {
				this.#refreshPromise = null
			})
		}

		return this.#refreshPromise
	}

	async #refreshAccessToken(): Promise<string> {
		while (true) {
			const { isAllowed } = await globalRequestLimiter.consume(donatexRateLimiterKey)
			if (!isAllowed) {
				await sleep(1000)
				continue
			}

			const response = await fetch('https://donatex.gg/api/connect/token', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/x-www-form-urlencoded',
				},
				body: new URLSearchParams({
					grant_type: 'refresh_token',
					refresh_token: this.refreshToken,
					client_id: config.DONATEX_CLIENT_ID!,
					client_secret: config.DONATEX_CLIENT_SECRET!,
				}).toString(),
			})

			if (!response.ok) {
				if (response.status === 429) {
					await sleep(1000)
					continue
				}

				throw new Error(
					`donatex token refresh failed with status ${response.status}: ${await response.text()}`
				)
			}

			const data = await response.json()
			this.accessToken = data.access_token
			if (data.refresh_token) {
				this.refreshToken = data.refresh_token
			}

			await updateDonateXIntegration({
				channel_id: this.twitchUserId,
				access_token: this.accessToken,
				refresh_token: this.refreshToken,
			})

			return this.accessToken
		}
	}

	async #donateCallback(data: DonateXDonation) {
		console.info(
			`[DONATEX #${this.twitchUserId}] Donation from ${data.username}: ${data.amount} ${data.currency}`
		)

		await onDonation({
			twitchUserId: this.twitchUserId,
			amount: data.amount,
			currency: data.currency,
			message: data.message,
			userName: data.username,
		})
	}

	async destroy() {
		if (this.#connection) {
			this.#connection.off('DonationCreated')
			await this.#connection.stop()
			this.#connection = null
		}
	}
}

export interface DonateXDonation {
	id: string
	username: string
	message: string | null
	currency: string
	amount: number
	amountInRub: number
	timestamp: string
	isTest: boolean
	isPotentiallyUnsafe: boolean
	isFeePaidByUser: boolean
	voiceFilePath?: string | null
}
