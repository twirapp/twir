import { createRouter, createWebHistory } from 'vue-router'

import { urqlClient } from './urql.js'

import type { RouteRecordRaw } from 'vue-router'

import { profileQuery, userAccessFlagChecker } from '@/api/auth.js'
import { ChannelRolePermissionEnum } from '@/gql/graphql.js'

export function newRouter() {
	const routes: ReadonlyArray<RouteRecordRaw> = [
		{
			path: '/dashboard/integrations/:integrationName',
			component: () => import('../pages/IntegrationsCallback.vue'),
		},
		{
			path: '/dashboard',
			component: () => import('../layout/layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('../pages/Dashboard.vue'),
					meta: {
						noPadding: true,
					},
				},
				{
					name: 'Integrations',
					path: '/dashboard/integrations',
					component: () => import('../pages/Integrations.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewIntegrations },
				},
				{
					path: '/dashboard/commands/:system',
					component: () => import('../features/commands/commands.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewCommands },
				},
				{
					path: '/dashboard/timers',
					component: () => import('../pages/Timers.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewTimers },
				},
				{
					path: '/dashboard/keywords',
					component: () => import('../pages/Keywords.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewKeywords },
				},
				{
					path: '/dashboard/variables',
					component: () => import('../pages/Variables.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewVariables },
				},
				{
					path: '/dashboard/greetings',
					component: () => import('../pages/greetings.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewGreetings,
					},
				},
				{
					path: '/dashboard/community',
					component: () => import('../pages/community.vue'),
					meta: { noPadding: true },
				},
				{
					path: '/dashboard/community/roles',
					component: () => import('../features/community-roles/community-roles.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewRoles, noPadding: true },
				},
				{
					path: '/dashboard/song-requests',
					component: () => import('../pages/SongRequests.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewSongRequests },
				},
				{
					path: '/dashboard/overlays',
					component: () => import('../pages/Overlays.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewOverlays },
				},
				{
					name: 'ChatOverlay',
					path: '/dashboard/overlays/chat',
					component: () => import('../pages/overlays/chat/Chat.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
						noPadding: true,
					},
				},
				{
					name: 'KappagenOverlay',
					path: '/dashboard/overlays/kappagen',
					component: () => import('../pages/overlays/kappagen/Kappagen.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
						noPadding: true,
					},
				},
				{
					name: 'BrbOverlay',
					path: '/dashboard/overlays/brb',
					component: () => import('../pages/overlays/brb/Brb.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
						noPadding: true,
					},
				},
				{
					name: 'DudesOverlay',
					path: '/dashboard/overlays/dudes',
					component: () => import('../pages/overlays/dudes/dudes-settings.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
						fullScreen: true,
					},
				},
				{
					path: '/dashboard/events/chat-alerts',
					component: () => import('../pages/chat-alerts.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewEvents, noPadding: true },
				},
				{
					path: '/dashboard/events/custom',
					component: () => import('../pages/Events.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewEvents },
				},
				{
					path: '/dashboard/alerts',
					component: () => import('../pages/alerts.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewAlerts },
				},
				{
					path: '/dashboard/games',
					component: () => import('../pages/Games.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewGames },
				},
				{
					path: '/dashboard/files',
					component: () => import('../pages/Files.vue'),
				},
				{
					name: 'RegistryOverlayEdit',
					path: '/dashboard/registry/overlays/:id',
					component: () => import('../components/registry/overlays/edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewOverlays },
				},
				{
					name: 'Moderation',
					path: '/dashboard/moderation',
					component: () => import('../pages/Moderation.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageModeration },
				},
				{
					name: 'Settings',
					path: '/dashboard/settings',
					component: () => import('../pages/user-settings/user-settings.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'AdminPanel',
					path: '/dashboard/admin',
					component: () => import('../pages/admin-panel.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'Import',
					path: '/dashboard/import',
					component: () => import('../pages/Import.vue'),
				},
				{
					name: 'Forbidden',
					path: '/dashboard/forbidden',
					component: () => import('../pages/NoAccess.vue'),
					meta: { fullScreen: true },
				},
				{
					path: '/:pathMatch(.*)*',
					name: 'NotFound',
					component: () => import('../pages/NotFound.vue'),
					meta: { fullScreen: true },
				},
			],
		},
	]

	const router = createRouter({
		history: createWebHistory(),
		routes,
	})

	router.beforeEach(async (to, _, next) => {
		try {
			const profileRequest = await urqlClient.value.executeQuery(profileQuery)
			if (!profileRequest.data) {
				return window.location.replace('/')
			}

			if (!to.meta.neededPermission) return next()

			const hasAccess = await userAccessFlagChecker(to.meta.neededPermission)
			if (hasAccess) {
				return next()
			}

			return next({ name: 'Forbidden' })
		} catch (error) {
			console.log(error)
			window.location.replace('/')
		}
	})

	return router
}
