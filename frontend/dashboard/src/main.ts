import './main.css';
import './assets/index.css';
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor';
import { broadcastQueryClient } from '@tanstack/query-broadcast-client-experimental';
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import urql, { cacheExchange, fetchExchange, subscriptionExchange } from '@urql/vue';
import { createClient as createWS, SubscribePayload } from 'graphql-ws';
import { createPinia } from 'pinia';
import { createApp } from 'vue';

import { i18n } from './plugins/i18n.js';
import { newRouter } from './plugins/router.js';

import App from '@/App.vue';

const pinia = createPinia();
const app = createApp(App);

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			refetchOnWindowFocus: false,
			refetchOnMount: false,
			refetchOnReconnect: false,
			retry: false,
		},
	},
});

broadcastQueryClient({
	queryClient,
	broadcastChannel: 'twir-dashboard',
});

VueQueryPlugin.install(app, {
	queryClient,
});

const meta = document.createElement('meta');
meta.name = 'naive-ui-style';
document.head.appendChild(meta);

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api-new/query`;
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api-new/query`;

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
});

app
	.use(pinia)
	.use(urql, {
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
	})
	.use(i18n)
	.use(newRouter(queryClient))
	.use(VueMonacoEditorPlugin);

app.mount('#app');

if (import.meta.env.DEV) {
	document.title = 'Twir (dev)';
}
