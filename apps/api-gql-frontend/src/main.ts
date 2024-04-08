import urql, { cacheExchange, fetchExchange } from '@urql/vue';
import { createApp } from 'vue';

import './style.css';
import App from './App.vue';

createApp(App).use(urql, {
	url: '/query',
	exchanges: [cacheExchange, fetchExchange],
	requestPolicy: 'cache-first',
	fetchOptions: {
		credentials: 'include',
	},
}).mount('#app');
