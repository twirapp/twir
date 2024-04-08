import urql, { cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue';
import { createClient as createWS } from 'graphql-ws';
import { createApp } from 'vue';

import './style.css';
import App from './App.vue';

const gqlWs = createWS({
  url: 'ws://localhost:3013/query',
  lazy: true,
});

createApp(App).use(urql, {
	url: '/query',
	exchanges: [
		cacheExchange,
		fetchExchange,
		subscriptionExchange({
			enableAllOperations: true,
      forwardSubscription: (operation) => ({
        subscribe: (sink) => ({
          unsubscribe: gqlWs.subscribe(operation, sink),
        }),
      }),
		}),
	],
	// requestPolicy: 'cache-first',
	fetchOptions: {
		credentials: 'include',
	},
}).mount('#app');
