import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { UsersGetResponse_UsersGetResponseUser as User, type UsersGetRequest, type UsersGetResponse } from '@twir/api/messages/admin_users/admin_users';
import type { Ref } from 'vue';

import { adminApiClient } from '../twirp';

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

	function mutateUser(userId: string, updater: (user: User) => Partial<User>) {
		queryClient.setQueriesData<UsersGetResponse>(
			['admin/users'],
			(oldData) => {
				if (!oldData) return oldData;

				return {
					total: oldData.total,
					users: oldData.users.map((user) => {
						if (user.id !== userId) return user;
						return {
							...user,
							...updater(user),
						};
					}),
				};
			},
		);
	}

	return {
		useUserSwitchBan: () => useMutation({
			mutationKey: ['admin/user/ban'],
			mutationFn: async (userId: string) => {
				const req = await adminApiClient.userSwitchBan({ userId });
				return req.response;
			},
			onSuccess: async (_, userId) => {
				mutateUser(userId, (user) => ({ isBanned: !user.isBanned }));
			},
		}),
		useUserSwitchAdmin: () => useMutation({
			mutationKey: ['admin/user/admin'],
			mutationFn: async (userId: string) => {
				const req = await adminApiClient.userSwitchAdmin({ userId });
				return req.response;
			},
			onSuccess: async (_, userId) => {
				mutateUser(userId, (user) => ({ isAdmin: !user.isAdmin }));
			},
		}),
	};
};
