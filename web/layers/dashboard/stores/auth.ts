import { createRequest } from '@urql/vue'
import { defineStore } from 'pinia'

import { graphql } from '~/gql/gql.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

export const dashboardUserProfile = createRequest(
	graphql(`
		query DashboardAuthenticatedUser {
			authenticatedUser {
				id
				isBotAdmin
				isBanned
				isEnabled
				isBotModerator
				hideOnLandingPage
				botId
				apiKey
				twitchProfile {
					description
					displayName
					login
					profileImageUrl
				}
				selectedDashboardId
				availableDashboards {
					id
					flags
					twitchProfile {
						login
						displayName
						profileImageUrl
					}
					apiKey
					plan {
						id
						name
						maxCommands
						maxTimers
						maxVariables
						maxAlerts
						maxEvents
						maxChatAlertsMessages
						maxCustomOverlays
						maxEightballAnswers
						maxCommandsResponses
						maxModerationRules
						maxKeywords
						maxGreetings
					}
				}
			}
		}
	`),
	{}
)

export const DashboardAuthStoreKey = 'dashboard-auth-user'

export const useDashboardAuth = defineStore('dashboard-auth', () => {
	const router = useRouter()

	const { data, executeQuery, fetching } = useQuery({
		query: dashboardUserProfile.query,
		context: {
			key: dashboardUserProfile.key,
			additionalTypenames: [DashboardAuthStoreKey],
		},
		variables: {},
	})

	const user = computed(() => {
		const authUser = data.value?.authenticatedUser
		if (!authUser) return null

		return {
			id: authUser.id,
			avatar: authUser.twitchProfile.profileImageUrl,
			login: authUser.twitchProfile.login,
			displayName: authUser.twitchProfile.displayName,
			apiKey: authUser.apiKey,
			isBotAdmin: authUser.isBotAdmin,
			isEnabled: authUser.isEnabled,
			isBanned: authUser.isBanned,
			isBotModerator: authUser.isBotModerator,
			botId: authUser.botId,
			selectedDashboardId: authUser.selectedDashboardId,
			hideOnLandingPage: authUser.hideOnLandingPage,
			availableDashboards: authUser.availableDashboards,
		}
	})

	// Track user with analytics (if rybbit is available)
	watch(
		user,
		(newUser) => {
			if (!newUser || !window?.rybbit || !import.meta.client) {
				return
			}

			window.rybbit.identify(newUser.id)
		},
		{ immediate: true }
	)

	const { executeMutation: executeLogout } = useMutation(graphql(`mutation userLogout { logout }`))

	async function logout() {
		await executeLogout({})
		if (typeof window !== 'undefined' && (window as any).rybbit) {
			;(window as any).rybbit.clearUserId()
		}
		await router.push('/')
	}

	function checkPermission(flag: ChannelRolePermissionEnum): boolean {
		if (!user.value) return false
		if (user.value.isBotAdmin) return true
		if (user.value.id === user.value.selectedDashboardId) return true

		const dashboard = user.value.availableDashboards.find(
			(d) => d.id === user.value?.selectedDashboardId
		)
		if (!dashboard) return false

		return (
			dashboard.flags.includes(ChannelRolePermissionEnum.CanAccessDashboard) ||
			dashboard.flags.includes(flag)
		)
	}

	return {
		user,
		isLoading: fetching,
		fetchUser: executeQuery,
		logout,
		checkPermission,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useDashboardAuth, import.meta.hot))
}

export const PERMISSIONS_FLAGS = [
	{ perm: ChannelRolePermissionEnum.CanAccessDashboard, description: 'All permissions' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.UpdateChannelTitle, description: 'Can update channel title' },
	{
		perm: ChannelRolePermissionEnum.UpdateChannelCategory,
		description: 'Can update channel category',
	},
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewCommands, description: 'Can view commands' },
	{ perm: ChannelRolePermissionEnum.ManageCommands, description: 'Can manage commands' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewKeywords, description: 'Can view keywords' },
	{ perm: ChannelRolePermissionEnum.ManageKeywords, description: 'Can manage keywords' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewTimers, description: 'Can view timers' },
	{ perm: ChannelRolePermissionEnum.ManageTimers, description: 'Can manage timers' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewIntegrations, description: 'Can view integrations' },
	{ perm: ChannelRolePermissionEnum.ManageIntegrations, description: 'Can manage integrations' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewSongRequests, description: 'Can view song requests' },
	{ perm: ChannelRolePermissionEnum.ManageSongRequests, description: 'Can manage song requests' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewModeration, description: 'Can view moderation settings' },
	{
		perm: ChannelRolePermissionEnum.ManageModeration,
		description: 'Can manage moderation settings',
	},
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewVariables, description: 'Can view variables' },
	{ perm: ChannelRolePermissionEnum.ManageVariables, description: 'Can manage variables' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewGreetings, description: 'Can view greetings' },
	{ perm: ChannelRolePermissionEnum.ManageGreetings, description: 'Can manage greetings' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewOverlays, description: 'Can view overlays' },
	{ perm: ChannelRolePermissionEnum.ManageOverlays, description: 'Can manage overlays' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewRoles, description: 'Can view roles' },
	{ perm: ChannelRolePermissionEnum.ManageRoles, description: 'Can manage roles' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewEvents, description: 'Can view events' },
	{ perm: ChannelRolePermissionEnum.ManageEvents, description: 'Can manage events' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewAlerts, description: 'Can view alerts' },
	{ perm: ChannelRolePermissionEnum.ManageAlerts, description: 'Can manage alerts' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewGames, description: 'Can view games' },
	{ perm: ChannelRolePermissionEnum.ManageGames, description: 'Can manage games' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewBotSettings, description: 'View bot settings' },
	{ perm: ChannelRolePermissionEnum.ManageBotSettings, description: 'Manage bot settings' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewModules, description: 'View modules' },
	{ perm: ChannelRolePermissionEnum.ManageModules, description: 'Manage modules' },
	'delimiter',
	{ perm: ChannelRolePermissionEnum.ViewGiveaways, description: 'View giveaways' },
	{ perm: ChannelRolePermissionEnum.ManageGiveaways, description: 'Manage giveaways' },
] as const

export function useUserAccessFlagChecker(flag: ChannelRolePermissionEnum) {
	const authStore = useDashboardAuth()

	return computed(() => authStore.checkPermission(flag))
}
