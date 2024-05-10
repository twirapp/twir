<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue'
import { type DropdownOption, NBadge, NButton, NDropdown } from 'naive-ui'
import { h } from 'vue'

import NotificationsDropdown from './notifications-dropdown.vue'

import { useNotifications } from '@/composables/use-notifications'

const { notificationsCounter } = useNotifications()

const profileOptions: DropdownOption[] = [
	{
		key: 'notifications',
		type: 'render',
		render: () => h(NotificationsDropdown),
	},
]
</script>

<template>
	<NDropdown
		trigger="click"
		:options="profileOptions"
		:on-update-show="notificationsCounter.onRead"
		size="large"
		style="top: 10px; left: 14px; width: 400px;"
	>
		<NBadge :value="notificationsCounter.counter">
			<NButton circle quaternary>
				<IconBell />
			</NButton>
		</NBadge>
	</NDropdown>
</template>
