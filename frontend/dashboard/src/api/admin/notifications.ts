import { useQuery, useSubscription } from '@urql/vue'
import { computed, watch } from 'vue'

import type { AdminNotificationsParams } from '@/gql/graphql'
import type { Ref } from 'vue'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

const invalidationKey = 'AdminNofiticationsInvalidateKey'

export function useQueryNotifications() {
	const { data: allNotifications } = useQuery({
		query: graphql(`
			query GetAllNotifications {
				notificationsByUser {
					id
					text
					createdAt
					editorJsJson
				}
			}
		`),
	})

	const { data: newNotifications } = useSubscription({
		query: graphql(`
			subscription NotificationsSubscription {
				newNotification {
					id
					text
					createdAt
					editorJsJson
				}
			}
		`),
	})

	watch(newNotifications, (newNotification) => {
		if (!allNotifications.value || !newNotification)
			return
		allNotifications.value.notificationsByUser.unshift(newNotification.newNotification)
	})

	const notifications = computed(() => {
		return allNotifications.value?.notificationsByUser ?? []
	})

	return notifications
}

export function useAdminNotifications() {
	const useQueryNotifications = (variables: Ref<AdminNotificationsParams>) => useQuery({
		context: {
			additionalTypenames: [invalidationKey],
		},
		get variables() {
			return {
				opts: variables.value,
			}
		},
		query: graphql(`
			query NotificationsByAdmin($opts: AdminNotificationsParams!) {
				notificationsByAdmin(opts: $opts) {
					total
					notifications {
						id
						text
						editorJsJson
						userId
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
						createdAt
					}
				}
			}
		`),
	})

	const useMutationCreateNotification = () => useMutation(graphql(`
		mutation CreateNotification($editorJsJson: String!, $userId: String) {
      notificationsCreate(editorJsJson: $editorJsJson, userId: $userId) {
				id
			}
    }
	`), [invalidationKey])

	const useMutationDeleteNotification = () => useMutation(graphql(`
		mutation DeleteNotification($id: ID!) {
			notificationsDelete(id: $id)
		}
	`), [invalidationKey])

	const useMutationUpdateNotifications = () => useMutation(graphql(`
		mutation UpdateNotifications($id: ID!, $opts: NotificationUpdateOpts!) {
			notificationsUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	return {
		useQueryNotifications,
		useMutationCreateNotification,
		useMutationDeleteNotification,
		useMutationUpdateNotifications,
	}
}
