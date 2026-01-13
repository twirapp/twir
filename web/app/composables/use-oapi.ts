import { Api, HttpClient } from '@twir/api/openapi'

export function useOapi(opts?: { headers?: Record<string, any> }) {
	const apiUrl = import.meta.server
		? import.meta.dev
			? 'http://localhost:3009'
			: 'http://api-gql:3009'
		: `${window.location.origin}/api`

	let headers = opts?.headers ?? {}
	try {
		if (import.meta.server && useRequestHeaders) {
			const serverHeaders = useRequestHeaders([
				'cookie',
				'session',
				'x-forwarded-for',
				'x-forwarded-proto',
				'x-forwarded-host',
				'x-real-ip',
				'cf-connecting-ip',
			])

			headers = {
				...headers,
				...serverHeaders,
			}
		}
	} catch {}

	return new Api(
		new HttpClient({
			baseUrl: apiUrl,
			baseApiParams: {
				credentials: 'include',
				headers,
			},
		})
	)
}
