import './main.css';
import './assets/index.css';
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor';
import { broadcastQueryClient } from '@tanstack/query-broadcast-client-experimental';
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import * as urql from '@urql/vue';
import { createPinia } from 'pinia';
import { createApp } from 'vue';

import { urqlClient } from './plugins/client.js';
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

app
	.use(pinia)
	.use(i18n)
	.use(urql, urqlClient)
	.use(newRouter(queryClient))
	.use(VueMonacoEditorPlugin)
	.mount('#app');

if (import.meta.env.DEV) {
	document.title = 'Twir (dev)';
}
