import { Api, HttpClient } from '@twir/api/openapi'

export interface TwirWebSocketEvent<T = Record<string, any>> {
	eventName: string
	data: T
	createdAt: string
}

export const openApi = new Api(
	new HttpClient({
		baseUrl: `${window.location.origin}/api`,
		baseApiParams: {
			credentials: 'include',
		},
	})
)
