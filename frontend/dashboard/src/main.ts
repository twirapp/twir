import './main.css';

import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor';
import { broadcastQueryClient } from '@tanstack/query-broadcast-client-experimental';
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { getProfile } from './api/index.js';
import { i18n } from './i18n.js';
import { newRouter } from './router.js';

import App from '@/App.vue';

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

getProfile(queryClient).catch(error => {
	console.error(error);
	window.location.replace('/');
}).then(() => {
	app
		.use(i18n)
		.use(newRouter(queryClient))
		.use(VueMonacoEditorPlugin);

	app.mount('#app');
});

