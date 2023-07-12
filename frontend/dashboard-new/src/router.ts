import { createRouter, createWebHistory } from 'vue-router';

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/',
			component: () => import('./pages/Home.vue'),
		},
		{
			path: '/integrations',
			component: () => import('./pages/Home.vue'),
		},
	],
});
