import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';

import { protectedClient, unprotectedClient } from '@/api/twirp';

export const useUserProfile = () => useQuery({
	queryKey: ['userProfile'],
	queryFn: async () => {
		const request = await protectedClient.authUserProfile({});
		return request.response;
	},
});

export const useLogout = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['userProfileLogout'],
		mutationFn: async () => {
			await protectedClient.authLogout({});
		},
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ['userProfile'],
			});
		},
	});
};

export const useLoginLink = () => useQuery({
	queryKey: ['loginLink'],
	queryFn: async () => {
		const redirectTo = window.location.href;
		const request = await unprotectedClient.authGetLink({ redirectTo });

		return request.response.link;
	},
});
