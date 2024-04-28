import { cacheExchange, Client, fetchExchange } from '@urql/vue';

const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api-new/query`;

export const urqlClient = new Client({
	url: gqlApiUrl,
	exchanges: [
		cacheExchange,
		fetchExchange,
	],
	// requestPolicy: 'cache-first',
	fetchOptions: {
		credentials: 'include',
	},
});
