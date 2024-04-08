import urql, { fetchExchange, cacheExchange } from '@urql/vue';
import { createApp } from 'vue';

import './style.css';
import App from './App.vue';

createApp(App).use(urql, {
	url: 'http://localhost:3002/query',
	exchanges: [cacheExchange, fetchExchange],
	requestPolicy: 'cache-first',
}).mount('#app');
