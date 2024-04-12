import { useQuery, useMutation, useSubscription } from '@urql/vue';
import { computed, watch, type Ref } from 'vue';

import { createCrudManager } from '../crud';
import { adminApiClient } from '../twirp';

import { graphql } from '@/gql';
import type { AdminNotificationsParams } from '@/gql/graphql';

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

export const useQueryNotifications = () => {
	const { data: allNotifications } = useQuery({
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

export const _useAdminNotifications = () => {
	const useQueryNotifications = (variables: Ref<AdminNotificationsParams>) => useQuery({
		get variables() {
			return {
				opts: variables.value,
			};
		},
		query: graphql(`
			query notificationsByAdmin($opts: AdminNotificationsParams!) {
				notificationsByAdmin(opts: $opts) {
					total
					notifications {
						id
						text
						userId
						twitchProfile {
							displayName
							profileImageUrl
						}
						createdAt
					}
				}
			}
		`),
	});

	const useMutationCreateNotification = () => useMutation(graphql(`
		mutation CreateNotification($text: String!, $userId: String) {
      notificationsCreate(text: $text, userId: $userId) {
				text
				userId
			}
    }
	`));

	return {
		useQueryNotifications,
		useMutationCreateNotification,
	};
};
