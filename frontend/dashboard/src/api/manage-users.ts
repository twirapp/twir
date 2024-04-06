import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { UsersGetRequest, UsersGetResponse } from '@twir/api/messages/admin_users/admin_users';
import type { Ref } from 'vue';

import { adminApiClient } from './twirp';

export const useAdminUsers = (payload: Ref<UsersGetRequest>) => useQuery({
	queryKey: ['admin/users', payload],
	keepPreviousData: true,
	queryFn: async () => {
		const req = await adminApiClient.getUsers(payload.value);
		return req.response;
	},
});

export const useAdminUserSwitcher = () => {
	const queryClient = useQueryClient();

	return {
		useUserSwitchBan: () => useMutation({
			mutationKey: ['admin/user/ban'],
			mutationFn: async (userId: string) => {
				const req = await adminApiClient.userSwitchBan({ userId });
				return req.response;
			},
			onSuccess: async (_, userId) => {
				queryClient.setQueriesData<UsersGetResponse>(
					['admin/users'],
					(oldData) => {
						if (!oldData) return oldData;

						return {
							...oldData,
							users: oldData.users.map((user) => {
								if (user.id !== userId) return user;
								return {
									...user,
									isBanned: !user.isBanned,
								};
							}),
						};
					},
				);
			},
		}),
		useUserSwitchAdmin: () => useMutation({
			mutationKey: ['admin/user/admin'],
			mutationFn: async (userId: string) => {
				const req = await adminApiClient.userSwitchAdmin({ userId });
				return req.response;
			},
			onSuccess: async (_, userId) => {
				queryClient.setQueriesData<UsersGetResponse>(
					['admin/users'],
					(oldData) => {
						if (!oldData) return oldData;

						return {
							...oldData,
							users: oldData.users.map((user) => {
								if (user.id !== userId) return user;
								return {
									...user,
									isAdmin: !user.isAdmin,
								};
							}),
						};
					},
				);
			},
		}),
	};
};
