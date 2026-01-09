import { useQueryNotifications } from '#layers/dashboard/api/admin/notifications'
import { createGlobalState, useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'

const NOTIFICATIONS_STORAGE_KEY = 'twirNotificationsCounter'

export const useNotifications = createGlobalState(() => {
	const notifications = useQueryNotifications()

	const notificationsStorage = useLocalStorage<string[]>(NOTIFICATIONS_STORAGE_KEY, [])

	const notificationsCounter = computed(() => {
		let notificationsCounter = 0

		const notificationsIds = notifications.value.map((notification) => notification.id)
		for (const notificationId of notificationsIds) {
			if (notificationsStorage.value.includes(notificationId)) continue
			notificationsCounter += 1
		}

		return {
			counter: notificationsCounter,
			onRead: () => (notificationsStorage.value = notificationsIds),
		}
	})

	return {
		notifications,
		notificationsCounter,
	}
})
