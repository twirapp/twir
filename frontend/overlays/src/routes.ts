import { createRouter, createWebHistory } from 'vue-router';

export const routes = createRouter({
	history: createWebHistory('/overlays'),
	routes: [
		{
			path: '/:apiKey/registry/overlays/:overlayId',
			component: () => import('@/pages/overlays.vue'),
		},
		{
			path: '/:apiKey/tts',
			component: () => import('@/pages/tts.vue'),
		},
		{
			path: '/:apiKey/obs',
			component: () => import('@/pages/obs.vue'),
		},
		{
			path: '/:apiKey/alerts',
			component: () => import('@/pages/alerts.vue'),
		},
		{
			path: '/:apiKey/chat',
			component: () => import('@/pages/overlays/chat.vue'),
		},
		{
			path: '/:apiKey/dudes',
			component: () => import('@/pages/overlays/dudes.vue'),
		},
		{
			path: '/:apiKey/kappagen',
			component: () => import('@/pages/overlays/kappagen.vue'),
		},
		{
			path: '/:apiKey/brb',
			component: () => import('@/pages/overlays/be-right-back.vue'),
		},
		{
			path: '/:apiKey/nowplaying',
			component: () => import('@/pages/overlays/nowplaying.vue'),
		},
	],
});
