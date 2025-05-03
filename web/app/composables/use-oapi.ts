import process from 'node:process'

import { Api, HttpClient } from '@twir/api/openapi'

import { useRequestHeaders } from '#imports'

export function useOapi() {
	const apiUrl = import.meta.server
		? process.env.NODE_ENV === 'production' ? 'http://api-gql:3009' : 'http://localhost:3009'
		: `${window.location.origin}/api`

	const headers = useRequestHeaders(['cookie', 'session'])

	return new Api(new HttpClient({
		baseUrl: apiUrl,
		baseApiParams: {
			credentials: 'include',
			headers,
		},
	}))
}
