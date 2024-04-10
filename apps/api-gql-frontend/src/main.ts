import urql, { cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue';
import { createClient as createWS, SubscribePayload } from 'graphql-ws';
import { createApp } from 'vue';

import './style.css';
import App from './App.vue';

const gqlWs = createWS({
  url: `ws://${window.location.host}/query`,
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
          unsubscribe: gqlWs.subscribe(operation as SubscribePayload, sink),
        }),
      }),
		}),
	],
	// requestPolicy: 'cache-first',
	fetchOptions: {
		credentials: 'include',
	},
}).mount('#app');
