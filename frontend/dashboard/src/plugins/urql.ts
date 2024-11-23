import { persistedExchange } from '@urql/exchange-persisted'
import { Client, cacheExchange, fetchExchange, mapExchange, subscriptionExchange } from '@urql/vue'
import { type SubscribePayload, createClient as createWS } from 'graphql-ws'
import { ref } from 'vue'

import { useToast } from '@/components/ui/toast'

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/query`
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api/query`

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
	shouldRetry: () => true,
})

const toast = useToast()

function createClient() {
	return new Client({
		url: gqlApiUrl,
		exchanges: [
			mapExchange({
				onError(error, _operation) {
					for (const er of error.graphQLErrors) {
						if (er.extensions.code !== 'BAD_REQUEST') continue
						const validationErrors = Object.entries(er.extensions.validation_errors ?? {} as Record<string, any>)

						toast.toast({
							variant: 'destructive',
							title: er.message,
							description: validationErrors.map(([key, value]) => `${key}: ${value}`).join(', '),
							duration: 10000,
						})
					}
				},
			}),
			cacheExchange,
			persistedExchange({
				preferGetForPersistedQueries: true,
				enableForMutation: true,
			}),
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
			const locationQuery = new URLSearchParams(window.location.search)
			const apiKey = locationQuery.get('apiKey')

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
