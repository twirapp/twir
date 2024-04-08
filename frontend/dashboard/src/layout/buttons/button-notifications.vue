<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue';
import { NButton, NBadge, NDropdown, type DropdownOption } from 'naive-ui';
import { h } from 'vue';
import { computed } from 'vue';

import NotificationsDropdown from './notifications-dropdown.vue';

import { useProtectedNotifications } from '@/api/admin/notifications';
import { useNotificationsCounter } from '@/composables/use-notifications-counter';

const { data } = useProtectedNotifications();
const { computeNotificationsCounter } = useNotificationsCounter();

const notifications = computed(() => {
	return data.value?.notifications ?? [];
});

const notificationsCounter = computed(() => {
	return computeNotificationsCounter(notifications.value);
});

function onUpdateShow(state: boolean) {
	if (state) return;
	notificationsCounter.value.readed();
}

const profileOptions: DropdownOption[] = [
	{
		key: 'notifications',
		type: 'render',
		render: () => h(NotificationsDropdown, { notifications: notifications.value }),
	},
];
</script>

<template>
	<n-dropdown
		trigger="click"
		:options="profileOptions"
		:on-update-show="onUpdateShow"
		size="large"
		style="top: 10px; left: 14px; width: 400px;"
	>
		<n-badge :value="notificationsCounter.counter">
			<n-button circle quaternary>
				<IconBell />
			</n-button>
		</n-badge>
	</n-dropdown>
</template>
