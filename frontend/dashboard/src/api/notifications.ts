import { useQuery } from '@tanstack/vue-query';

import { createCrudManager } from './crud';
import { adminApiClient, protectedApiClient } from './twirp';

export const useAdminNotifications = () => createCrudManager({
	queryKey: 'admin/notifications',
	create: adminApiClient.notificationsCreate,
	deleteOne: adminApiClient.notificationsDelete,
	update: adminApiClient.notificationsUpdate,
	getAll: adminApiClient.notificationsGetAll,
	getOne: null,
	patch: null,
});

export const useProtectedNotifications = () => useQuery({
	queryKey: ['protected/notification'],
	queryFn: async () => {
		const req = await protectedApiClient.notificationsGetAll({});
		return req.response;
	},
});
