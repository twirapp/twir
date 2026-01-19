import type { SSRExchange } from '@urql/vue';
import { cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue';
import type { SubscribePayload } from 'graphql-ws';
import { createClient as createWS } from 'graphql-ws';

import { defineUrqlClient } from '#urql/client';

export default defineUrqlClient((ssrExchange) => {
	const exchanges = import.meta.server ? setupServer(ssrExchange) : setupClient(ssrExchange)

	const headers = useRequestHeaders([
		'cookie',
		'session',
		'x-forwarded-for',
		'x-forwarded-proto',
		'x-forwarded-host',
		'x-real-ip',
		'cf-connecting-ip',
		'X-Ru-Detected-IP',
		'remote-host',
	])

	return {
		exchanges,
		fetchOptions: {
			credentials: 'include',
			headers,
		},
	}
})

function setupServer(ssrExchange: SSRExchange) {
	return [cacheExchange, ssrExchange, fetchExchange]
}

function setupClient(ssrExchange: SSRExchange) {
	const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/query`
	const gqlWs = createWS({
		url: wsUrl!,
		lazy: true,
		shouldRetry: () => true,
	})

	return [
		cacheExchange,
		ssrExchange,
		fetchExchange,
		subscriptionExchange({
			enableAllOperations: true,
			forwardSubscription: (operation) => ({
				subscribe: (sink) => ({
					unsubscribe: gqlWs!.subscribe(operation as SubscribePayload, sink),
				}),
			}),
		}),
	]
}
