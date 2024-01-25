import { storeToRefs } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';

import { useStreamerProfile } from '@/composables/use-streamer-profile';

export const publicRouteNames = {
	commands: 'Commands',
	songRequests: 'Song Requests',
	ttsSettings: 'TTS Settings',
	users: 'Users',
};

const routeNames = {
	notFound: 'Not found',
};

export const createPublicRouter = () => {
	const profileStore = useStreamerProfile();
	const { profile } = storeToRefs(profileStore);

	const router = createRouter({
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
				children: [
					{
						name: publicRouteNames.commands,
						path: '/p/:channelName',
						alias: '/p/:channelName/commands',
						component: () => import('./pages/commands.vue'),
					},
					{
						name: publicRouteNames.songRequests,
						path: '/p/:channelName/songs-requests',
						component: () => import('./pages/song-requests.vue'),
					},
					{
						name: publicRouteNames.ttsSettings,
						path: '/p/:channelName/tts-settings',
						component: () => import('./pages/tts.vue'),
					},
					{
						name: publicRouteNames.users,
						path: '/p/:channelName/users',
						component: () => import('./pages/users.vue'),
					},
				],
			},
		],
	});

	router.beforeEach(async (to) => {
		if (to.name === routeNames.notFound) return true;

		const channelName = to.params.channelName;
		if (typeof channelName !== 'string') {
			return {
				name: routeNames.notFound,
			};
		}

		profileStore.fetchProfile(channelName).finally(() => {
			if (!profile.value) {
				router.push({ name: routeNames.notFound });
			}
		});
	});

	return router;
};


