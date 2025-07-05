import { cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue'
import { createClient as createWS } from 'graphql-ws'

import type { ClientOptions } from '@urql/vue'
import type { SubscribePayload } from 'graphql-ws'

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/query`
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api/query`

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
	shouldRetry: () => true,
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
}
