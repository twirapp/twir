import {
	QueryClient,
	QueryOptions,
	useMutation,
	useQuery,
	useQueryClient,
} from '@tanstack/vue-query';
import type { Dashboard } from '@twir/api/messages/auth/auth';
import { useQuery as useGraphqlQuery } from '@urql/vue';
import { computed } from 'vue';

import { protectedApiClient } from './twirp.js';

import { graphql } from '@/gql';

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

export const dashboardsQueryOptions = {
	queryKey: ['authDashboards'],
	queryFn: async () => {
		const call = await protectedApiClient.authGetDashboards({});
		return call.response;
	},
} satisfies QueryOptions;

export const useDashboards = () => useQuery(dashboardsQueryOptions);

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

export const PERMISSIONS_FLAGS = {
	CAN_ACCESS_DASHBOARD: 'All permissions',
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
	const dashboards = useDashboards();

	return computed(() => {
		if (!dashboards.data.value?.dashboards || !profile.data.value?.selectedDashboardId) return false;

		if (profile.data.value.id == profile.data.value.selectedDashboardId) {
			return true;
		}

		const dashboard = dashboards.data.value.dashboards.find(d => d.id === profile.data.value!.selectedDashboardId);
		if (!dashboard) return false;

		if (dashboard.flags.includes('CAN_ACCESS_DASHBOARD')) return true;
		return dashboard.flags.includes(flag);
	});
};

export const userAccessFlagChecker = async (queryClient: QueryClient, flag: PermissionsType) => {
	const { data: profile, executeQuery: fetchProfile } = useProfile();
	await fetchProfile();

	const { dashboards } = await queryClient.getQueryData(dashboardsQueryOptions.queryKey) as {
		dashboards: Dashboard[]
	};

	if (!dashboards || !profile.value || !profile?.value?.selectedDashboardId) return false;
	if (profile.value.selectedDashboardId == profile.value.id) return true;

	const dashboard = dashboards.find(d => d.id === profile.value!.selectedDashboardId);
	if (!dashboard) return false;

	if (dashboard.flags.includes('CAN_ACCESS_DASHBOARD')) return true;
	return dashboard.flags.includes(flag);
};
