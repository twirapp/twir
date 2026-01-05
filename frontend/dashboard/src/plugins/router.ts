import { type RouteRecordRaw, createRouter, createWebHistory } from 'vue-router'

import { profileQuery, userAccessFlagChecker } from '@/api/auth.js'
import { ChannelRolePermissionEnum } from '@/gql/graphql.js'

import { urqlClient } from './urql.js'

export function newRouter() {
	const routes: ReadonlyArray<RouteRecordRaw> = [
		{
			path: '/dashboard/integrations/spotify',
			component: () => import('@/pages/IntegrationsCallbackSpotify.vue'),
		},
		{
			path: '/dashboard/integrations/donationalerts',
			component: () => import('@/pages/IntegrationsCallbackDonationAlerts.vue'),
		},
		{
			path: '/dashboard/integrations/nightbot',
			component: () => import('@/features/import/nightbot/nightbot-callback.vue'),
		},
		{
			path: '/dashboard/integrations/valorant',
			component: () => import('@/features/integrations/pages/valorant-callback.vue'),
		},
		{
			path: '/dashboard/integrations/discord',
			component: () => import('@/features/integrations/pages/discord-callback.vue'),
		},
		{
			path: '/dashboard/integrations/vk',
			component: () => import('@/features/integrations/pages/vk-callback.vue'),
		},
		{
			path: '/dashboard/integrations/:integrationName',
			component: () => import('@/pages/IntegrationsCallback.vue'),
		},
		{
			path: '/dashboard/popup',
			component: () => import('@/popup-layout/popup-layout.vue'),
			children: [
				{
					path: '/dashboard/popup/widgets/eventslist',
					component: () => import('@/components/dashboard/events.vue'),
					props: { popup: true },
				},
				{
					path: '/dashboard/popup/widgets/audit-log',
					component: () => import('@/components/dashboard/audit-logs.vue'),
					props: { popup: true },
				},
			],
		},
		{
			path: '/dashboard',
			component: () => import('@/layout/layout.vue'),
			children: [
				{
					path: '/dashboard',
					component: () => import('@/pages/Dashboard.vue'),
					meta: { noPadding: true },
				},
				{
					path: '/dashboard/bot-settings',
					component: () => import('@/features/bot-settings/bot-settings.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ViewBotSettings,
						noPadding: true,
					},
				},
				{
					name: 'Integrations',
					path: '/dashboard/integrations',
					component: () => import('@/pages/Integrations.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewIntegrations, noPadding: true },
				},
				{
					name: 'DiscordIntegration',
					path: '/dashboard/integrations/discord-settings',
					component: () => import('@/features/integrations/pages/discord-settings.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewIntegrations, noPadding: true },
				},
				{
					path: '/dashboard/commands/:system',
					component: () => import('@/features/commands/commands.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewCommands, noPadding: true },
				},
				{
					path: '/dashboard/commands/:system/:id',
					component: () => import('@/features/commands/commands-edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageCommands, noPadding: true },
				},
				{
					path: '/dashboard/timers',
					component: () => import('@/features/timers/timers.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewTimers, noPadding: true },
				},
				{
					path: '/dashboard/giveaways',
					component: () => import('@/features/giveaways/giveaways.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewGiveaways, noPadding: true },
				},
				{
					name: 'giveaways-view',
					path: '/dashboard/giveaways/view/:id',
					component: () => import('@/features/giveaways/giveaways.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewGiveaways, noPadding: true },
				},
				{
					path: '/dashboard/modules',
					component: () => import('@/features/modules/modules.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewModules, noPadding: true },
				},
				{
					path: '/dashboard/timers/:id',
					component: () => import('@/features/timers/timers-edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageTimers, noPadding: true },
				},
				{
					path: '/dashboard/keywords',
					component: () => import('@/pages/Keywords.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewKeywords, noPadding: true },
				},
				{
					path: '/dashboard/variables',
					component: () => import('@/features/variables/variables.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewVariables, noPadding: true },
				},
				{
					path: '/dashboard/variables/:id',
					component: () => import('@/features/variables/variables-edit.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageVariables, noPadding: true },
				},
				{
					path: '/dashboard/greetings',
					component: () => import('@/pages/greetings.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewGreetings,
					},
				},
				{
					path: '/dashboard/expiring-vips',
					component: () => import('@/features/expiring-vips/expiring-vips.vue'),
					meta: { noPadding: true },
				},
				{
					path: '/dashboard/community',
					component: () => import('@/pages/community.vue'),
					meta: { noPadding: true },
				},
				{
					path: '/dashboard/community/roles',
					component: () => import('@/features/community-roles/community-roles.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ViewRoles,
						noPadding: true,
					},
				},
				{
					path: '/dashboard/song-requests',
					component: () => import('@/pages/SongRequests.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewSongRequests },
				},
				{
					path: '/dashboard/overlays',
					component: () => import('@/pages/Overlays.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewOverlays },
				},
				{
					name: 'ChatOverlay',
					path: '/dashboard/overlays/chat',
					component: () => import('@/pages/overlays/chat/Chat.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'KappagenOverlay',
					path: '/dashboard/overlays/kappagen',
					component: () => import('@/features/overlays/kappagen/kappagen.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'BrbOverlay',
					path: '/dashboard/overlays/brb',
					component: () => import('@/features/overlays/brb/page.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'TTSOverlay',
					path: '/dashboard/overlays/tts',
					component: () => import('@/features/overlays/tts/page.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'ObsOverlay',
					path: '/dashboard/overlays/obs',
					component: () => import('@/features/overlays/obs/page.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'DudesOverlay',
					path: '/dashboard/overlays/dudes',
					component: () => import('@/pages/overlays/dudes/dudes-settings.vue'),
					meta: {
						fullScreen: false,
						neededPermission: ChannelRolePermissionEnum.ManageOverlays,
					},
				},
				{
					name: 'FaceitStatsOverlay',
					path: '/dashboard/overlays/faceit-stats',
					component: () => import('@/features/overlays/faceit-stats/builder.vue'),
					meta: {
						noPadding: true,
					},
				},
				{
					name: 'ValorantStatsOverlay',
					path: '/dashboard/overlays/valorant-stats',
					component: () => import('@/features/overlays/valorant-stats/builder.vue'),
					meta: {
						noPadding: true,
					},
				},
				{
					path: '/dashboard/events/chat-alerts',
					component: () => import('@/pages/chat-alerts.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewEvents,
					},
				},
				{
					path: '/dashboard/events',
					component: () => import('@/features/events/events-list.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewEvents, noPadding: true },
				},
				{
					path: '/dashboard/events/:id',
					component: () => import('@/features/events/event-form.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageEvents, noPadding: true },
				},
				{
					path: '/dashboard/alerts',
					component: () => import('@/pages/alerts.vue'),
					meta: {
						noPadding: true,
						neededPermission: ChannelRolePermissionEnum.ViewAlerts,
					},
				},
				{
					path: '/dashboard/games',
					component: () => import('@/pages/Games.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ViewGames, noPadding: true },
				},
				{
					path: '/dashboard/files',
					component: () => import('@/pages/Files.vue'),
				},
				{
					name: 'RegistryOverlayEdit',
					path: '/dashboard/registry/overlays/:id',
					component: () => import('@/components/registry/overlays/edit.vue'),
					meta: {
						neededPermission: ChannelRolePermissionEnum.ViewOverlays,
						fullScreen: true,
					},
				},
				{
					name: 'Moderation',
					path: '/dashboard/moderation',
					component: () => import('@/features/moderation/moderation.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageModeration, noPadding: true },
				},
				{
					name: 'ModerationForm',
					path: '/dashboard/moderation/:id',
					component: () => import('@/features/moderation/moderation-form.vue'),
					meta: { neededPermission: ChannelRolePermissionEnum.ManageModeration, noPadding: true },
				},
				{
					name: 'Settings',
					path: '/dashboard/settings',
					component: () => import('@/pages/user-settings/user-settings.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'AdminPanel',
					path: '/dashboard/admin',
					component: () => import('@/pages/admin-panel.vue'),
					meta: { noPadding: true, adminOnly: true },
				},
				{
					name: 'Import',
					path: '/dashboard/import',
					component: () => import('@/pages/Import.vue'),
				},
				{
					name: 'Notifications',
					path: '/dashboard/notifications',
					component: () => import('@/pages/notifications.vue'),
					meta: { noPadding: true },
				},
				{
					name: 'Forbidden',
					path: '/dashboard/forbidden',
					component: () => import('@/pages/NoAccess.vue'),
					meta: { fullScreen: true },
				},
				{
					path: '/:pathMatch(.*)*',
					name: 'NotFound',
					component: () => import('@/pages/NotFound.vue'),
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
		if (to.path.startsWith('/dashboard/popup')) return next()

		try {
			const profileRequest = await urqlClient.value.executeQuery(profileQuery)
			if (!profileRequest.data) {
				return window.location.replace('/')
			}

			if (to.meta.adminOnly && !profileRequest.data.authenticatedUser.isBotAdmin) {
				return next({ name: 'NotFound' })
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
