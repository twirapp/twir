import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import MainApp from './app.vue';
import { router } from './routes.js';

import './style.css';

const app = createApp(MainApp);
app.use(router);

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
