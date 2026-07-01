import { config } from '@twir/config'

import { validateUrl, type ValidationResult } from './url-validation'

const FETCH_TIMEOUT_MS = 10000
const MAX_BODY_SIZE = 1_048_576 // 1MB
const isProd = config.NODE_ENV === 'production'
const PROXY_URL = 'http://warp-proxy:9091'

export interface FetchResponse {
	status: number
	statusText: string
	headers: Record<string, string>
	body: string
}

export async function hostFetch(url: string, options?: RequestInit): Promise<FetchResponse> {
	const validation: ValidationResult = await validateUrl(url)

	const controller = new AbortController()
	const timeout = setTimeout(() => controller.abort(), FETCH_TIMEOUT_MS)

	const fetchUrl = rewriteUrl(url, validation.resolvedIp)

	console.log(`[hostFetch] ${options?.method || 'GET'} ${url} (proxy: ${isProd})`)

	try {
		const headers = new Headers(options?.headers)
		if (validation.resolvedIp) {
			const parsed = new URL(url)
			headers.set('Host', parsed.host)
		}

		const response = await fetch(fetchUrl, {
			...options,
			headers,
			signal: controller.signal,
			redirect: 'error',
			...(isProd ? { proxy: PROXY_URL } : {}),
		})

		const responseHeaders: Record<string, string> = {}
		response.headers.forEach((value, key) => {
			responseHeaders[key] = value
		})

		const body = await readBodyWithLimit(response)

		console.log(`[hostFetch] Done: status=${response.status}, body length=${body.length}`)

		return {
			status: response.status,
			statusText: response.statusText,
			headers: responseHeaders,
			body,
		}
	} catch (err: any) {
		console.error(`[hostFetch] Error:`, err.message)
		throw err
	} finally {
		clearTimeout(timeout)
	}
}

function rewriteUrl(url: string, resolvedIp: string | null): string {
	if (!resolvedIp) return url

	const parsed = new URL(url)
	const isIpv6 = resolvedIp.includes(':')
	parsed.hostname = isIpv6 ? `[${resolvedIp}]` : resolvedIp
	return parsed.toString()
}

async function readBodyWithLimit(response: Response): Promise<string> {
	const reader = response.body?.getReader()
	if (!reader) return ''

	const chunks: Uint8Array[] = []
	let totalSize = 0

	try {
		while (true) {
			const { done, value } = await reader.read()
			if (done) break
			totalSize += value.length
			if (totalSize > MAX_BODY_SIZE) {
				reader.cancel()
				throw new Error(`Response body exceeded ${MAX_BODY_SIZE} bytes`)
			}
			chunks.push(value)
		}
	} catch (err) {
		reader.cancel()
		throw err
	}

	return new TextDecoder().decode(Buffer.concat(chunks))
}
