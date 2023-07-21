import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { Profile } from '@twir/grpc/generated/api/api/auth';
import { ref } from 'vue';

import { protectedApiClient } from './twirp.js';

const profile = ref<Profile | null>(null);


const profileQueryKey = ['authProfile'];
export const useProfile = () =>
	useQuery<Profile | null>({
		queryKey: profileQueryKey,
		queryFn: async () => {
			const user = await getProfile();
			profile.value = user;
			return user;
		},
		retry: false,
		initialData: profile.value,
	});

export const useLogout = () => useMutation({
	mutationKey: ['authLogout'],
	mutationFn: async () => {
		await protectedApiClient.authLogout({});
		profile.value = null;
	},
	onSuccess: () => {
		profile.value = null;
		window.location.replace('/');
	},
});

export const getProfile = async () => {
	if (profile.value) return profile.value;

	const call = await protectedApiClient.authUserProfile({});
	profile.value = call.response;
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
			await queryClient.invalidateQueries(dashboardsQueryKey);
			await queryClient.invalidateQueries(profileQueryKey);
		},
	});
};
