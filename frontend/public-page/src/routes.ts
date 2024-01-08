import { createRouter, createWebHistory } from 'vue-router';

export const routesNames = {
	commands: 'Commands',
	songRequests: 'Song Requests',
	ttsSettings: 'TTS Settings',
	users: 'Users',
};

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/p/:channelName',
			component: () => import('./components/layout.vue'),
			children: [
				{
					name: routesNames.commands,
					path: '/p/:channelName',
					alias: '/p/:channelName/commands',
					component: () => import('./pages/commands.vue'),
				},
				{
					name: routesNames.songRequests,
					path: '/p/:channelName/songs-requests',
					component: () => import('./pages/song-requests.vue'),
				},
				{
					name: routesNames.ttsSettings,
					path: '/p/:channelName/tts-settings',
					component: () => import('./pages/tts-settings.vue'),
				},
				{
					name: routesNames.users,
					path: '/p/:channelName/users',
					component: () => import('./pages/users.vue'),
				},
			],
		},
	],
});
