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
					name: 'Integrations',
					path: '/dashboard/integrations',
					component: () => import('./pages/Integrations.vue'),
				},
				{
					path: '/dashboard/integrations/:integrationName',
					component: () => import('./pages/IntegrationsCallback.vue'),
				},
				{
					path: '/dashboard/commands/:system',
					component: () => import('./pages/Commands.vue'),
				},
				{
					path: '/dashboard/timers',
					component: () => import('./pages/Timers.vue'),
				},
				{
					path: '/dashboard/keywords',
					component: () => import('./pages/Keywords.vue'),
				},
				{
					path: '/dashboard/variables',
					component: () => import('./pages/Variables.vue'),
				},
				{
					path: '/dashboard/greetings',
					component: () => import('./pages/Greetings.vue'),
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
