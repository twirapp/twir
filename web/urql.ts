import type { SSRExchange } from '@urql/vue'
import type { SubscribePayload } from 'graphql-ws'

import { defineUrqlClient } from '#urql/client'
import { cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue'
import { createClient as createWS } from 'graphql-ws'

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
		preferGetMethod: false,
		fetchOptions: {
			credentials: 'include',
			headers: {
				...headers,
				...getApiKeyHeader(),
			},
		},
	}
})

function getApiKeyHeader(): Record<string, string> {
	if (typeof window === 'undefined') return {}

	// Check query param first (legacy)
	const params = new URLSearchParams(window.location.search)
	const queryKey = params.get('apiKey')
	if (queryKey) return { 'Api-Key': queryKey }

	// Check widget route: /w/{channelApiKey}/...
	const widgetMatch = window.location.pathname.match(/^\/w\/([^/]+)\//)
	if (widgetMatch?.[1]) return { 'Api-Key': widgetMatch[1] }

	// Check overlay route: /o/{apiKey}/...
	const overlayMatch = window.location.pathname.match(/^\/o\/([^/]+)\//)
	if (overlayMatch?.[1]) return { 'Api-Key': overlayMatch[1] }

	return {}
}

function setupServer(ssrExchange: SSRExchange) {
	return [cacheExchange, ssrExchange, fetchExchange]
}

function setupClient(ssrExchange: SSRExchange) {
	const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/query`
	let acknowledgedConnections = 0
	const gqlWs = createWS({
		url: wsUrl!,
		lazy: true,
		retryAttempts: Infinity,
		shouldRetry: () => true,
		on: {
			connected: () => {
				acknowledgedConnections += 1
				if (acknowledgedConnections > 1) {
					window.dispatchEvent(new CustomEvent('twir:graphql-ws-reconnected'))
				}
			},
		},
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
