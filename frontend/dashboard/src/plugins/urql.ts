import { cacheExchange, Client, fetchExchange, subscriptionExchange } from '@urql/vue';
import { createClient as createWS, type SubscribePayload } from 'graphql-ws';
import { ref } from 'vue';

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api-new/query`;
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api-new/query`;

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
	shouldRetry: () => true,
});

export const urqlClient = ref<Client | null>(null);

export const useUrqlClient = () => {
	function createClient() {
		return new Client({
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
			// requestPolicy: 'cache-first',
			fetchOptions: {
				credentials: 'include',
			},
		});
	}

	return {
		urqlClient,
		createClient
	}
}
