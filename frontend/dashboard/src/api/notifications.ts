import { useQuery } from '@tanstack/vue-query';

import { createCrudManager } from './crud';
import { adminApiClient, protectedApiClient } from './twirp';

export const useAdminNotifications = () => createCrudManager({
	client: adminApiClient,
	queryKey: 'admin/notifications',
	create: adminApiClient.notificationsCreate,
	deleteOne: adminApiClient.notificationsDelete,
	update: adminApiClient.notificationsUpdate,
	getAll: adminApiClient.notificationsGetAll,
	getOne: null,
	patch: null,
	invalidateAdditionalQueries: ['protected/notifications'],
});

export const useProtectedNotifications = () => useQuery({
	queryKey: ['protected/notifications'],
	queryFn: async () => {
		const req = await protectedApiClient.notificationsGetAll({});
		return req.response;
	},
});
