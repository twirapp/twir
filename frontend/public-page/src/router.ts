import { IconCommand, IconHeadphones, IconUsers } from '@tabler/icons-vue';
import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';

import IconTts from '@/assets/icons/sidebar/tts.svg?use';

export const publicRouteNames = {
	commands: 'Commands',
	songRequests: 'Song Requests',
	ttsSettings: 'TTS Settings',
	users: 'Users',
};

export const routeNames = {
	notFound: 'Not found',
};

export const channelRoutes: RouteRecordRaw[] = [
	{
		name: publicRouteNames.commands,
		path: '/p/:channelName',
		alias: '/p/:channelName/commands',
		component: () => import('./pages/commands.vue'),
		meta: {
			icon: IconCommand,
		},
	},
	{
		name: publicRouteNames.songRequests,
		path: '/p/:channelName/songs-requests',
		component: () => import('./pages/song-requests.vue'),
		meta: {
			icon: IconHeadphones,
		},
	},
	{
		name: publicRouteNames.ttsSettings,
		path: '/p/:channelName/tts-settings',
		component: () => import('./pages/tts.vue'),
		meta: {
			icon: IconTts,
		},
	},
	{
		name: publicRouteNames.users,
		path: '/p/:channelName/users',
		component: () => import('./pages/users.vue'),
		meta: {
			icon: IconUsers,
		},
	},
];

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/p/404',
			name: routeNames.notFound,
			component: () => import('./pages/404.vue'),
		},
		{
			path: '/p/:channelName',
			component: () => import('./layout/layout.vue'),
			children: channelRoutes,
		},
	],
});


