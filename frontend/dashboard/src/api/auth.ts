import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { Profile } from '@twir/grpc/generated/api/api/auth';
import { computed } from 'vue';

import { protectedApiClient } from './twirp.js';


const profileQueryKey = ['authProfile'];
export const useProfile = () =>
	useQuery<Profile | null>({
		queryKey: profileQueryKey,
		queryFn: async () => {
			const call = await protectedApiClient.authUserProfile({});
			return call.response;
		},
		retry: false,
	});

export const useLogout = () => useMutation({
	mutationKey: ['authLogout'],
	mutationFn: async () => {
		await protectedApiClient.authLogout({});
	},
	onSuccess: () => {
		window.location.replace('/');
	},
});

export const getProfile = async () => {
	const call = await protectedApiClient.authUserProfile({});
	return call.response;
};

const dashboardsQueryKey = ['authDashboards'];

export const useDashboards = () => useQuery({
	queryKey: dashboardsQueryKey,
	queryFn: async () => {
		const call = await protectedApiClient.authGetDashboards({});
		return call.response;
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
		},
	});
};


export const useUserAccessFlagChecker = (flag: string) => {
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
