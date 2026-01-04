import { type ClientOptions, cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue'
import { type SubscribePayload, createClient as createWS } from 'graphql-ws'
import { ref } from 'vue'

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/query`
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api/query`

// Store API key in a ref that can be accessed outside component context
export const apiKeyRef = ref<string | null>(null)

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
	shouldRetry: () => true,
	connectionParams: () => {
		if (apiKeyRef.value) {
			return {
				'api-key': apiKeyRef.value,
			}
		}
		return {}
	},
})

export const urqlClientOptions: ClientOptions = {
	url: gqlApiUrl,
	exchanges: [
		cacheExchange,
		fetchExchange,
		subscriptionExchange({
			enableAllOperations: true,
			forwardSubscription: (operation) => ({
				subscribe: (sink) => ({
					unsubscribe: gqlWs.subscribe(operation as SubscribePayload, sink),
				}),
			}),
		}),
	],
	fetchOptions: () => {
		const headers: Record<string, string> = {}

		if (apiKeyRef.value) {
			headers['api-key'] = apiKeyRef.value
			headers['x-api-key'] = apiKeyRef.value
		}

		return {
			headers,
			credentials: 'include',
		}
	},
}
