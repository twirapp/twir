import { Client, cacheExchange, fetchExchange } from '@urql/vue'

const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api/query`

export const urqlClient = new Client({
	url: gqlApiUrl,
	exchanges: [
		cacheExchange,
		fetchExchange
	],
	// requestPolicy: 'cache-first',
	fetchOptions: {
		credentials: 'include'
	}
})
