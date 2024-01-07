import { createRouter, createWebHistory } from 'vue-router';

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/p/:channelName',
			component: () => import('./components/layout.vue'),
			children: [
				{
					name: 'Commands',
					path: '/p/:channelName',
					alias: '/p/:channelName/commands',
					component: () => import('./pages/commands.vue'),
				},
				{
					name: 'Song requests',
					path: '/p/:channelName/songs-requests',
					component: () => import('./pages/song-requests.vue'),
				},
				{
					name: 'TTS Settings',
					path: '/p/:channelName/tts-settings',
					component: () => import('./pages/tts-settings.vue'),
				},
				{
					name: 'Users',
					path: '/p/:channelName/users',
					component: () => import('./pages/users.vue'),
				},
			],
		},
	],
});
