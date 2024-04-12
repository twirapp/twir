<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue';
import { NButton, NBadge, NDropdown, type DropdownOption } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { h } from 'vue';

import NotificationsDropdown from './notifications-dropdown.vue';

import { useNotifications } from '@/composables/use-notifications';

const { notificationsCounter } = storeToRefs(useNotifications());

const profileOptions: DropdownOption[] = [
	{
		key: 'notifications',
		type: 'render',
		render: () => h(NotificationsDropdown),
	},
];
</script>

<template>
	<n-dropdown
		trigger="click"
		:options="profileOptions"
		:on-update-show="notificationsCounter.onRead"
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
