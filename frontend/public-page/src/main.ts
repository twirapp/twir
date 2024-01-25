import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import { createPinia } from 'pinia';
import { createApp } from 'vue';

import MainApp from './app.vue';
import { createPublicRouter } from './routes.js';

import './assets/index.css';

const pinia = createPinia();
const app = createApp(MainApp)
	.use(pinia)
	.use(createPublicRouter());

VueQueryPlugin.install(app, {
	queryClient: new QueryClient({
		defaultOptions: {
			queries: {
				refetchOnWindowFocus: false,
				refetchOnMount: false,
				refetchOnReconnect: false,
				retry: false,
			},
		},
	}),
});

app.mount('#app');
