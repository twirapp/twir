import { useQuery } from '@tanstack/vue-query';
import { useQuery as _useQuery, useSubscription } from '@urql/vue';
import { computed } from 'vue';

import { createCrudManager } from '../crud';
import { adminApiClient, protectedApiClient } from '../twirp';

import { graphql } from '@/gql';
import type { NotificationsGetAllQuery } from '@/gql/graphql';

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

	const notifications = computed<NotificationsGetAllQuery['notificationsByUser']>(() => {
		const notification = [...(allNotifications.value?.notificationsByUser ?? [])];

		if (newNotifications.value) {
			notification.push(newNotifications.value.newNotification);
		}

		return notification;
	});

	return notifications;
};
