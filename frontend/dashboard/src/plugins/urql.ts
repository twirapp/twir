import { Client, cacheExchange, fetchExchange, mapExchange, subscriptionExchange } from '@urql/vue'
import { type SubscribePayload, createClient as createWS } from 'graphql-ws'
import { ref } from 'vue'
import { toast } from 'vue-sonner'

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/query`
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api/query`

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
	shouldRetry: () => true,
	connectionParams() {
		const apiKey = getApiKeyFromUrlQuery()

		if (apiKey) {
			return {
				'api-key': apiKey,
			}
		}

		return {}
	},
})

function createClient() {
	return new Client({
		url: gqlApiUrl,
		exchanges: [
			mapExchange({
				onError(error, _operation) {
					for (const er of error.graphQLErrors) {
						if (er.extensions.code !== 'BAD_REQUEST') continue
						const validationErrors = Object.entries(
							er.extensions.validation_errors ?? ({} as Record<string, any>)
						)

						toast.error(er.message, {
							description: validationErrors.map(([key, value]) => `${key}: ${value}`).join(', '),
							duration: 10000,
						})
					}
				},
			}),
			cacheExchange,
			// persistedExchange({
			// 	preferGetForPersistedQueries: true,
			// 	enableForMutation: true,
			// 	generateHash: (_, document) => {
			// 		// eslint-disable-next-line ts/ban-ts-comment
			// 		// @ts-expect-error
			// 		return document.__meta__.hash
			// 	},
			// }),
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
		// requestPolicy: 'cache-first',
		fetchOptions: () => {
			const apiKey = getApiKeyFromUrlQuery()

			const options: RequestInit = {
				credentials: 'include',
			}

			if (apiKey) {
				options.headers = {
					...options.headers,
					'Api-Key': apiKey,
				}
			}

			return options
		},
	})
}

export const urqlClient = ref<Client>(createClient())

export function useUrqlClient() {
	function reInitClient() {
		urqlClient.value = createClient()
	}

	return {
		urqlClient,
		reInitClient,
	}
}

function getApiKeyFromUrlQuery() {
	const locationQuery = new URLSearchParams(window.location.search)
	const apiKey = locationQuery.get('apiKey')

	return apiKey
}
