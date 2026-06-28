import { config } from '@twir/config'

import { validateUrl } from './url-validation'

const FETCH_TIMEOUT_MS = 10000
const isProd = config.NODE_ENV === 'production'
const PROXY_URL = 'http://warp-proxy:9091'

export interface FetchResponse {
	status: number
	statusText: string
	headers: Record<string, string>
	body: string
}

export async function hostFetch(url: string, options?: RequestInit): Promise<FetchResponse> {
	await validateUrl(url)

	const controller = new AbortController()
	const timeout = setTimeout(() => controller.abort(), FETCH_TIMEOUT_MS)

	try {
		const response = await fetch(url, {
			...options,
			signal: controller.signal,
			...(isProd ? { proxy: PROXY_URL } : {}),
		})

		const headers: Record<string, string> = {}
		response.headers.forEach((value, key) => {
			headers[key] = value
		})

		const body = await response.text()

		return {
			status: response.status,
			statusText: response.statusText,
			headers,
			body,
		}
	} finally {
		clearTimeout(timeout)
	}
}
