import { QueryClient } from '@tanstack/vue-query';
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

import { type PermissionsType, userAccessFlagChecker } from '@/api/index.js';

type Route = Omit<RouteRecordRaw, 'meta' | 'children'> & {
	meta?: { neededPermission?: PermissionsType; noPadding?: boolean },
	children?: ReadonlyArray<Route>
}

export const newRouter = (queryClient: QueryClient) => {
	const routes: ReadonlyArray<Route> = [
		{
			path: '/dashboard',
			component: () => import('./layout/layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('./pages/Dashboard.vue'),
					meta: {
						noPadding: true,
					},
				},
				{
					name: 'Integrations',
					path: '/dashboard/integrations',
					component: () => import('./pages/Integrations.vue'),
					meta: { neededPermission: 'VIEW_INTEGRATIONS' },
				},
				{
					path: '/dashboard/integrations/:integrationName',
					component: () => import('./pages/IntegrationsCallback.vue'),
				},
				{
					path: '/dashboard/commands/:system',
					component: () => import('./pages/Commands.vue'),
					meta: { neededPermission: 'VIEW_COMMANDS' },
				},
				{
					path: '/dashboard/timers',
					component: () => import('./pages/Timers.vue'),
					meta: { neededPermission: 'VIEW_TIMERS' },
				},
				{
					path: '/dashboard/keywords',
					component: () => import('./pages/Keywords.vue'),
					meta: { neededPermission: 'VIEW_KEYWORDS' },
				},
				{
					path: '/dashboard/variables',
					component: () => import('./pages/Variables.vue'),
					meta: { neededPermission: 'VIEW_VARIABLES' },
				},
				{
					path: '/dashboard/greetings',
					component: () => import('./pages/Greetings.vue'),
					meta: { neededPermission: 'VIEW_GREETINGS' },
				},
				{
					path: '/dashboard/community/users',
					component: () => import('./pages/CommunityUsers.vue'),
				},
				{
					path: '/dashboard/community/roles',
					component: () => import('./pages/CommunityRoles.vue'),
					meta: { neededPermission: 'VIEW_ROLES' },
				},
				{
					path: '/dashboard/song-requests',
					component: () => import('./pages/SongRequests.vue'),
					meta: { neededPermission: 'VIEW_SONG_REQUESTS' },
				},
				{
					path: '/dashboard/overlays',
					component: () => import('./pages/Overlays.vue'),
					meta: { neededPermission: 'VIEW_OVERLAYS' },
				},
				{
					path: '/dashboard/events',
					component: () => import('./pages/Events.vue'),
					meta: { neededPermission: 'VIEW_EVENTS' },
				},
				{
					path: '/dashboard/alerts',
					component: () => import('./pages/Alerts.vue'),
					meta: { neededPermission: 'VIEW_ALERTS' },
				},
				{
					path: '/dashboard/files',
					component: () => import('./pages/Files.vue'),
					// meta: { neededPermission: 'VIEW_EVENTS' },
				},
				{
					name: 'Forbidden',
					path: '/dashboard/forbidden',
					component: () => import('./pages/NoAccess.vue'),
				},
			],
		},
	];

	const router = createRouter({
		history: createWebHistory(),
		// eslint-disable-next-line @typescript-eslint/ban-ts-comment
		// @ts-ignore
		routes,
	});

	router.beforeEach(async (to, _, next) => {
		if (!to.meta.neededPermission) return next();

		const hasAccess = await userAccessFlagChecker(queryClient, to.meta.neededPermission as PermissionsType);
		if (hasAccess) {
			return next();
		}

		return next({ name: 'Forbidden' });
	});

	return router;
};
