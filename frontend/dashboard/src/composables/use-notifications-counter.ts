import type { Notification } from '@twir/api/messages/notifications/notifications';
import { useLocalStorage } from '@vueuse/core';
import { defineStore } from 'pinia';

export const useNotificationsCounter = defineStore('notifications-counter', () => {
	const readedNotifications = useLocalStorage<string[]>('twirNotificationsCounter', []);

	function computeNotificationsCounter(notifications: Notification[]) {
		const notif: string[] = [];

		for (const notification of notifications) {
			if (readedNotifications.value.includes(notification.id)) continue;
			notif.push(notification.id);
		}

		return {
			counter: notif.length,
			readed: () => readedNotifications.value.push(...notif),
		};
	}

	return {
		computeNotificationsCounter,
	};
});
