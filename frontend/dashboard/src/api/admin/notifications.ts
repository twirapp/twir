import { useQuery, useMutation, useSubscription } from '@urql/vue';
import { computed, watch, type Ref } from 'vue';

import { graphql } from '@/gql';
import type { AdminNotificationsParams } from '@/gql/graphql';

export const useQueryNotifications = () => {
	const { data: allNotifications } = useQuery({
		query: graphql(`
			query GetAllNotifications {
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
		allNotifications.value.notificationsByUser.unshift(newNotification.newNotification);
	});

	const notifications = computed(() => {
		return allNotifications.value?.notificationsByUser ?? [];
	});

	return notifications;
};

export const useAdminNotifications = () => {
	const useQueryNotifications = (variables: Ref<AdminNotificationsParams>) => useQuery({
		get variables() {
			return {
				opts: variables.value,
			};
		},
		query: graphql(`
			query NotificationsByAdmin($opts: AdminNotificationsParams!) {
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
				id
			}
    }
	`));

	const useMutationDeleteNotification = () => useMutation(graphql(`
		mutation DeleteNotification($id: ID!) {
			notificationsDelete(id: $id)
		}
	`));

	const useMutationUpdateNotifications = () => useMutation(graphql(`
		mutation UpdateNotifications($id: ID!, $opts: NotificationUpdateOpts!) {
			notificationsUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`));

	return {
		useQueryNotifications,
		useMutationCreateNotification,
		useMutationDeleteNotification,
		useMutationUpdateNotifications,
	};
};
