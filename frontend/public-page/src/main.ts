import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';

import './style.css';
import App from './App.vue';

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/p/:channelName',
			component: () => import('./layout/Layout.vue'),
			children: [
				{
					name: 'Commands',
					path: '/p/:channelName',
					alias: '/p/:channelName/commands',
					component: () => import('./pages/Commands.vue'),
				},
				{
					name: 'Song requests',
					path: '/p/:channelName/songs-requests',
					component: () => import('./pages/SongRequests.vue'),
				},
				{
					name: 'TTS Settings',
					path: '/p/:channelName/tts-settings',
					component: () => import('./pages/TTSSettings.vue'),
				},
			],
		},
	],
});

const app = createApp(App)
	.use(router);

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
