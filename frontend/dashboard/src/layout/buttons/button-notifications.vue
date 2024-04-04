<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue';
import { NButton, NBadge, NDropdown, type DropdownOption } from 'naive-ui';
import { h } from 'vue';
import { computed } from 'vue';

import NotificationsDropdown from './notifications-dropdown.vue';

import { useProtectedNotifications } from '@/api/notifications';

const { data } = useProtectedNotifications();

const notifications = computed(() => {
	if (!data.value) return [];
	return data.value.notifications;
});

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
		size="large"
		style="top: 10px; left: 14px; width: 500px;"
	>
		<n-badge :value="notifications.length">
			<n-button circle quaternary>
				<IconBell />
			</n-button>
		</n-badge>
	</n-dropdown>
</template>
