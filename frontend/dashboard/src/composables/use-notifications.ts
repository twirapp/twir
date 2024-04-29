import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import { computed } from 'vue'

import { useQueryNotifications } from '@/api/admin/notifications.js'

const NOTIFICATIONS_STORAGE_KEY = 'twirNotificationsCounter'

export const useNotifications = defineStore('notifications', () => {
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
			onRead: (state: boolean) => {
				if (state) return
				notificationsStorage.value = notificationsIds
			}
		}
	})

	return {
		notifications,
		notificationsCounter
	}
})
