import { useMutation, useQueryClient } from '@tanstack/vue-query';
import { useQuery as useGraphqlQuery } from '@urql/vue';
import { computed } from 'vue';

import { protectedApiClient } from './twirp.js';

import { graphql } from '@/gql';
import { ChannelRolePermissionEnum } from '@/gql/graphql.js';
import { urqlClient } from '@/plugins/urql';

export const useProfileQuery = graphql(`
	query AuthenticatedUser {
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
			}
		}
	}
`);

export const useProfile = () => {
	const { data: response, executeQuery, fetching } = useGraphqlQuery({
		query: useProfileQuery,
	});

	const computedUser = computed(() => {
		const user = response.value?.authenticatedUser;
		if (!user) return null;

		return {
			id: user.id,
			avatar: user.twitchProfile.profileImageUrl,
			login: user.twitchProfile.login,
			displayName: user.twitchProfile.displayName,
			apiKey: user.apiKey,
			isBotAdmin: user.isBotAdmin,
			isEnabled: user.isEnabled,
			isBanned: user.isBanned,
			isBotModerator: user.isBotModerator,
			botId: user.botId,
			selectedDashboardId: user.selectedDashboardId,
			hideOnLandingPage: user.hideOnLandingPage,
			availableDashboards: user.availableDashboards,
		};
	});

	return { data: computedUser, executeQuery, isLoading: fetching };
};

export const useLogout = () => useMutation({
	mutationKey: ['authLogout'],
	mutationFn: async () => {
		await protectedApiClient.authLogout({});
	},
});

export const useSetDashboard = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['authSetDashboard'],
		mutationFn: async (dashboardId: string) => {
			await protectedApiClient.authSetDashboard({ dashboardId });
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries();
			await queryClient.resetQueries();
		},
	});
};

export const PERMISSIONS_FLAGS: { [key in ChannelRolePermissionEnum]: string } = {
	CAN_ACCESS_DASHBOARD: 'All permissions',
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-expect-error
	empty1: '',

	UPDATE_CHANNEL_TITLE: 'Can update channel title',
	UPDATE_CHANNEL_CATEGORY: 'Can update channel category',

	VIEW_COMMANDS: 'Can view commands',
	MANAGE_COMMANDS: 'Can manage commands',

	VIEW_KEYWORDS: 'Can view keywords',
	MANAGE_KEYWORDS: 'Can manage keywords',

	VIEW_TIMERS: 'Can view timers',
	MANAGE_TIMERS: 'Can manage timers',

	VIEW_INTEGRATIONS: 'Can view integrations',
	MANAGE_INTEGRATIONS: 'Can manage integrations',

	VIEW_SONG_REQUESTS: 'Can view song requests',
	MANAGE_SONG_REQUESTS: 'Can manage song requests',

	VIEW_MODERATION: 'Can view moderation settings',
	MANAGE_MODERATION: 'Can manage moderation settings',

	VIEW_VARIABLES: 'Can view variables',
	MANAGE_VARIABLES: 'Can manage variables',

	VIEW_GREETINGS: 'Can view greetings',
	MANAGE_GREETINGS: 'Can manage greetings',

	VIEW_OVERLAYS: 'Can view overlays',
	MANAGE_OVERLAYS: 'Can manage overlays',

	VIEW_ROLES: 'Can view roles',
	MANAGE_ROLES: 'Can manage roles',

	VIEW_EVENTS: 'Can view events',
	MANAGE_EVENTS: 'Can manage events',

	VIEW_ALERTS: 'Can view alerts',
	MANAGE_ALERTS: 'Can manage alerts',

	VIEW_GAMES: 'Can view games',
	MANAGE_GAMES: 'Can manage games',
};

export type PermissionsType = keyof typeof PERMISSIONS_FLAGS

export const useUserAccessFlagChecker = (flag: PermissionsType) => {
	const profile = useProfile();

	return computed(() => {
		if (!profile.data.value?.availableDashboards || !profile.data.value?.selectedDashboardId) return false;

		if (profile.data.value.id == profile.data.value.selectedDashboardId) {
			return true;
		}

		if (profile.data.value.isBotAdmin) return true;

		const dashboard = profile.data.value?.availableDashboards.find(dashboard => {
			return dashboard.id === profile.data.value!.selectedDashboardId;
		});
		if (!dashboard) return false;

		if (dashboard.flags.includes(ChannelRolePermissionEnum.CanAccessDashboard)) return true;
		return dashboard.flags.includes(flag as ChannelRolePermissionEnum);
	});
};

export const userAccessFlagChecker = async (flag: PermissionsType) => {
	const { data: profile } = await urqlClient.executeQuery({ query: useProfileQuery, key: 0, variables: {} });

	if (profile?.authenticatedUser.isBotAdmin) return true;
	if (!profile || !profile?.authenticatedUser.selectedDashboardId) return false;
	if (profile.authenticatedUser.selectedDashboardId == profile.authenticatedUser.id) return true;

	const dashboard = profile.authenticatedUser.availableDashboards.find(dashboard => {
		return dashboard.id === profile.authenticatedUser.selectedDashboardId;
	});
	if (!dashboard) return false;

	if (dashboard.flags.includes(ChannelRolePermissionEnum.CanAccessDashboard)) return true;
	return dashboard.flags.includes(flag as ChannelRolePermissionEnum);
};
