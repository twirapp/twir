import { useQuery } from '@tanstack/vue-query';
import type { UsersGetRequest } from '@twir/api/messages/admin_users/admin_users';
import type { Ref } from 'vue';

import { adminApiClient } from './twirp'	;

export const useAdminUsers = (payload: Ref<UsersGetRequest>) => useQuery({
	queryKey: ['admin/users', payload],
	queryFn: async () => {
		const req = await adminApiClient.getUsers(payload.value);
		return req.response;
	},
});
