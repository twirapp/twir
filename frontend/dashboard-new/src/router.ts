import { createRouter, createWebHistory } from 'vue-router';

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/dashboard',
			component: () => import('./pages/Home.vue'),
		},
		{
			path: '/dashboard/integrations',
			component: () => import('./pages/Home.vue'),
		},
	],
});
