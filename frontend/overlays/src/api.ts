import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'
import { UnProtectedClient } from '@twir/api/api.client'
import { Api, HttpClient } from '@twir/api/openapi'

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api-old/v1`,
	sendJson: import.meta.env.DEV,
})

export const unprotectedApiClient = new UnProtectedClient(transport)

export interface TwirWebSocketEvent<T = Record<string, any>> {
	eventName: string
	data: T
	createdAt: string
}

export const openApi = new Api(new HttpClient({
	baseUrl: `${window.location.origin}/api`,
	baseApiParams: {
		credentials: 'include',
	},
}))
