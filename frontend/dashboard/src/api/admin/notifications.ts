import { useQuery } from '@tanstack/vue-query';
import { useQuery as _useQuery, useSubscription } from '@urql/vue';
import { computed, watch } from 'vue';

import { createCrudManager } from '../crud';
import { adminApiClient, protectedApiClient } from '../twirp';

import { graphql } from '@/gql';

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

export const useQueryNotifications = () => {
	const { data: allNotifications } = _useQuery({
		query: graphql(`
			query NotificationsGetAll {
				notificationsByUser {
					id
					text
					createdAt
				}
			}
		`),
	});

	const { data: newNotifications } = useSubscription({
		query: graphql(`
			subscription NotificationsSubscription {
				newNotification {
					id
					text
					createdAt
				}
			}
		`),
	});

	watch(newNotifications, (newNotification) => {
		if (!allNotifications.value || !newNotification) return;
		allNotifications.value.notificationsByUser.push(newNotification.newNotification);
	});

	const notifications = computed(() => {
		return allNotifications.value?.notificationsByUser ?? [];
	});

	return notifications;
};
