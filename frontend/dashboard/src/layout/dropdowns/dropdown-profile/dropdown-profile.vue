<script setup lang="ts">
import { type DropdownOption, NAvatar, NButton, NDropdown, NSpin } from 'naive-ui'
import { h } from 'vue'

import DropdownFooter from './dropdown-profile-footer.vue'
import DropdownHeader from './dropdown-profile-header.vue'

import { useProfile } from '@/api/auth.js'

const { data: profileData, isLoading: isProfileLoading } = useProfile()

// TODO: Close dropdown when clicking to options
const profileOptions: DropdownOption[] = [
	{
		key: 'header',
		type: 'render',
		render: () => h(DropdownHeader),
	},
	{
		key: 'header-divider',
		type: 'divider',
	},
	{
		key: 'footer',
		type: 'render',
		render: () => h(DropdownFooter),
	},
]
</script>

<template>
	<NDropdown
		trigger="click"
		:options="profileOptions"
		size="large"
		style="top: 12px; left: 14px; width: 400px;"
	>
		<NButton text>
			<NSpin v-if="isProfileLoading" size="small" />

			<div v-else class="flex gap-1 items-center">
				<NAvatar
					size="small"
					:src="profileData?.avatar"
					round
				/>
			</div>
		</NButton>
	</NDropdown>
</template>
