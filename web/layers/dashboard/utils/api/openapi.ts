import { Api, HttpClient } from '@twir/api/openapi'

export const openApi = new Api(new HttpClient({
	baseUrl: `${window.location.origin}/api`,
	baseApiParams: {
		credentials: 'include',
	},
}))
