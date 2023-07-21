import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { Profile } from '@twir/grpc/generated/api/api/auth';

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
			// await queryClient.invalidateQueries(dashboardsQueryKey);
			// await queryClient.invalidateQueries(profileQueryKey);
		},
	});
};
