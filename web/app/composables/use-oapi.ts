import process from 'node:process'

import { Api, HttpClient } from '@twir/api/openapi'

export function useOapi(opts?: { headers?: Record<string, any> }) {
	const apiUrl = import.meta.server
		? process.env.NODE_ENV === 'production' ? 'http://api-gql:3009' : 'http://localhost:3009'
		: `${window.location.origin}/api`

	let headers = {}
	try {
		if (opts?.headers) {
			headers = opts.headers
		} else if (import.meta.server && useRequestHeaders) {
			headers = useRequestHeaders(['cookie', 'session'])
		}
	} catch {}

	return new Api(new HttpClient({
		baseUrl: apiUrl,
		baseApiParams: {
			credentials: 'include',
			headers,
		},
	}))
}
