import { createRouter, createWebHistory } from 'vue-router';

import { getProfile } from './api/index.js';

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/dashboard',
			component: () => import('./layout/layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('./pages/Dashboard.vue'),
				},
				{
					path: '/dashboard/integrations',
					component: () => import('./pages/Dashboard.vue'),
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
