import { createRouter, createWebHistory } from 'vue-router';

import { getProfile } from './api/index.js';

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/dashboard',
			component: () => import('./Layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('./pages/Home.vue'),
				},
				{
					path: '/dashboard/integrations',
					component: () => import('./pages/Home.vue'),
				},
			],
		},
	],
});

router.beforeEach(async () => {
	try {
		await getProfile();
		return true;
	} catch (e) {
		console.error(e);
		window.location.replace('/');
		return false;
	}
});
